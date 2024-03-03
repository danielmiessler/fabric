import os
import io
import json
import httpx
import binascii
import platform
import urllib.parse
from os import PathLike
from pathlib import Path
from hashlib import sha256
from base64 import b64encode, b64decode

from typing import Any, AnyStr, Union, Optional, Sequence, Mapping, Literal,TypedDict

import sys

if sys.version_info < (3, 9):
  from typing import Iterator, AsyncIterator
else:
  from collections.abc import Iterator, AsyncIterator


if sys.version_info < (3, 11):
  from typing_extensions import NotRequired
else:
  from typing import NotRequired


from importlib import metadata



class BaseGenerateResponse(TypedDict):
  model: str
  'Model used to generate response.'

  created_at: str
  'Time when the request was created.'

  done: bool
  'True if response is complete, otherwise False. Useful for streaming to detect the final response.'

  total_duration: int
  'Total duration in nanoseconds.'

  load_duration: int
  'Load duration in nanoseconds.'

  prompt_eval_count: int
  'Number of tokens evaluated in the prompt.'

  prompt_eval_duration: int
  'Duration of evaluating the prompt in nanoseconds.'

  eval_count: int
  'Number of tokens evaluated in inference.'

  eval_duration: int
  'Duration of evaluating inference in nanoseconds.'


class GenerateResponse(BaseGenerateResponse):
  """
  Response returned by generate requests.
  """

  response: str
  'Response content. When streaming, this contains a fragment of the response.'

  context: Sequence[int]
  'Tokenized history up to the point of the response.'


class Message(TypedDict):
  """
  Chat message.
  """

  role: Literal['user', 'assistant', 'system']
  "Assumed role of the message. Response messages always has role 'assistant'."

  content: str
  'Content of the message. Response messages contains message fragments when streaming.'

  images: NotRequired[Sequence[Any]]
  """
  Optional list of image data for multimodal models.

  Valid input types are:

  - `str` or path-like object: path to image file
  - `bytes` or bytes-like object: raw image data

  Valid image formats depend on the model. See the model card for more information.
  """


class ChatResponse(BaseGenerateResponse):
  """
  Response returned by chat requests.
  """

  message: Message
  'Response message.'


class ProgressResponse(TypedDict):
  status: str
  completed: int
  total: int
  digest: str


class Options(TypedDict, total=False):
  # load time options
  numa: bool
  num_ctx: int
  num_batch: int
  num_gqa: int
  num_gpu: int
  main_gpu: int
  low_vram: bool
  f16_kv: bool
  logits_all: bool
  vocab_only: bool
  use_mmap: bool
  use_mlock: bool
  embedding_only: bool
  rope_frequency_base: float
  rope_frequency_scale: float
  num_thread: int

  # runtime options
  num_keep: int
  seed: int
  num_predict: int
  top_k: int
  top_p: float
  tfs_z: float
  typical_p: float
  repeat_last_n: int
  temperature: float
  repeat_penalty: float
  presence_penalty: float
  frequency_penalty: float
  mirostat: int
  mirostat_tau: float
  mirostat_eta: float
  penalize_newline: bool
  stop: Sequence[str]


class RequestError(Exception):
  """
  Common class for request errors.
  """

  def __init__(self, error: str):
    super().__init__(error)
    self.error = error
    'Reason for the error.'


class ResponseError(Exception):
  """
  Common class for response errors.
  """

  def __init__(self, error: str, status_code: int = -1):
    try:
      # try to parse content as JSON and extract 'error'
      # fallback to raw content if JSON parsing fails
      error = json.loads(error).get('error', error)
    except json.JSONDecodeError:
      ...

    super().__init__(error)
    self.error = error
    'Reason for the error.'

    self.status_code = status_code
    'HTTP status code of the response.'


class BaseClient:
  def __init__(
    self,
    client,
    host: Optional[str] = None,
    follow_redirects: bool = True,
    timeout: Any = None,
    **kwargs,
  ) -> None:
    """
    Creates a httpx client. Default parameters are the same as those defined in httpx
    except for the following:
    - `follow_redirects`: True
    - `timeout`: None
    `kwargs` are passed to the httpx client.
    """

    headers = kwargs.pop('headers', {})
    headers['Content-Type'] = 'application/json'
    headers['Accept'] = 'application/json'
    headers['User-Agent'] = f'({platform.machine()} {platform.system().lower()}) Python/{platform.python_version()}'

    self._client = client(
      base_url=_parse_host(host or os.getenv('OLLAMA_HOST')),
      follow_redirects=follow_redirects,
      timeout=timeout,
      headers=headers,
      **kwargs,
    )


class Client(BaseClient):
  def __init__(self, host: Optional[str] = None, **kwargs) -> None:
    super().__init__(httpx.Client, host, **kwargs)

  def _request(self, method: str, url: str, **kwargs) -> httpx.Response:
    response = self._client.request(method, url, **kwargs)

    try:
      response.raise_for_status()
    except httpx.HTTPStatusError as e:
      raise ResponseError(e.response.text, e.response.status_code) from None

    return response

  def _stream(self, method: str, url: str, **kwargs) -> Iterator[Mapping[str, Any]]:
    with self._client.stream(method, url, **kwargs) as r:
      try:
        r.raise_for_status()
      except httpx.HTTPStatusError as e:
        e.response.read()
        raise ResponseError(e.response.text, e.response.status_code) from None

      for line in r.iter_lines():
        partial = json.loads(line)
        if e := partial.get('error'):
          raise ResponseError(e)
        yield partial

  def _request_stream(
    self,
    *args,
    stream: bool = False,
    **kwargs,
  ) -> Union[Mapping[str, Any], Iterator[Mapping[str, Any]]]:
    return self._stream(*args, **kwargs) if stream else self._request(*args, **kwargs).json()

  def generate(
    self,
    model: str = '',
    prompt: str = '',
    system: str = '',
    template: str = '',
    context: Optional[Sequence[int]] = None,
    stream: bool = False,
    raw: bool = False,
    format: Literal['', 'json'] = '',
    images: Optional[Sequence[AnyStr]] = None,
    options: Optional[Options] = None,
    keep_alive: Optional[Union[float, str]] = None,
  ) -> Union[Mapping[str, Any], Iterator[Mapping[str, Any]]]:
    """
    Create a response using the requested model.

    Raises `RequestError` if a model is not provided.

    Raises `ResponseError` if the request could not be fulfilled.

    Returns `GenerateResponse` if `stream` is `False`, otherwise returns a `GenerateResponse` generator.
    """

    if not model:
      raise RequestError('must provide a model')

    return self._request_stream(
      'POST',
      '/api/generate',
      json={
        'model': model,
        'prompt': prompt,
        'system': system,
        'template': template,
        'context': context or [],
        'stream': stream,
        'raw': raw,
        'images': [_encode_image(image) for image in images or []],
        'format': format,
        'options': options or {},
        'keep_alive': keep_alive,
      },
      stream=stream,
    )

  def chat(
    self,
    model: str = '',
    messages: Optional[Sequence[Message]] = None,
    stream: bool = False,
    format: Literal['', 'json'] = '',
    options: Optional[Options] = None,
    keep_alive: Optional[Union[float, str]] = None,
  ) -> Union[Mapping[str, Any], Iterator[Mapping[str, Any]]]:
    """
    Create a chat response using the requested model.

    Raises `RequestError` if a model is not provided.

    Raises `ResponseError` if the request could not be fulfilled.

    Returns `ChatResponse` if `stream` is `False`, otherwise returns a `ChatResponse` generator.
    """

    if not model:
      raise RequestError('must provide a model')

    for message in messages or []:
      if not isinstance(message, dict):
        raise TypeError('messages must be a list of Message or dict-like objects')
      if not (role := message.get('role')) or role not in ['system', 'user', 'assistant']:
        raise RequestError('messages must contain a role and it must be one of "system", "user", or "assistant"')
      if not message.get('content'):
        raise RequestError('messages must contain content')
      if images := message.get('images'):
        message['images'] = [_encode_image(image) for image in images]

    return self._request_stream(
      'POST',
      '/api/chat',
      json={
        'model': model,
        'messages': messages,
        'stream': stream,
        'format': format,
        'options': options or {},
        'keep_alive': keep_alive,
      },
      stream=stream,
    )

  def embeddings(
    self,
    model: str = '',
    prompt: str = '',
    options: Optional[Options] = None,
    keep_alive: Optional[Union[float, str]] = None,
  ) -> Sequence[float]:
    return self._request(
      'POST',
      '/api/embeddings',
      json={
        'model': model,
        'prompt': prompt,
        'options': options or {},
        'keep_alive': keep_alive,
      },
    ).json()

  def pull(
    self,
    model: str,
    insecure: bool = False,
    stream: bool = False,
  ) -> Union[Mapping[str, Any], Iterator[Mapping[str, Any]]]:
    """
    Raises `ResponseError` if the request could not be fulfilled.

    Returns `ProgressResponse` if `stream` is `False`, otherwise returns a `ProgressResponse` generator.
    """
    return self._request_stream(
      'POST',
      '/api/pull',
      json={
        'name': model,
        'insecure': insecure,
        'stream': stream,
      },
      stream=stream,
    )

  def push(
    self,
    model: str,
    insecure: bool = False,
    stream: bool = False,
  ) -> Union[Mapping[str, Any], Iterator[Mapping[str, Any]]]:
    """
    Raises `ResponseError` if the request could not be fulfilled.

    Returns `ProgressResponse` if `stream` is `False`, otherwise returns a `ProgressResponse` generator.
    """
    return self._request_stream(
      'POST',
      '/api/push',
      json={
        'name': model,
        'insecure': insecure,
        'stream': stream,
      },
      stream=stream,
    )

  def create(
    self,
    model: str,
    path: Optional[Union[str, PathLike]] = None,
    modelfile: Optional[str] = None,
    stream: bool = False,
  ) -> Union[Mapping[str, Any], Iterator[Mapping[str, Any]]]:
    """
    Raises `ResponseError` if the request could not be fulfilled.

    Returns `ProgressResponse` if `stream` is `False`, otherwise returns a `ProgressResponse` generator.
    """
    if (realpath := _as_path(path)) and realpath.exists():
      modelfile = self._parse_modelfile(realpath.read_text(), base=realpath.parent)
    elif modelfile:
      modelfile = self._parse_modelfile(modelfile)
    else:
      raise RequestError('must provide either path or modelfile')

    return self._request_stream(
      'POST',
      '/api/create',
      json={
        'name': model,
        'modelfile': modelfile,
        'stream': stream,
      },
      stream=stream,
    )

  def _parse_modelfile(self, modelfile: str, base: Optional[Path] = None) -> str:
    base = Path.cwd() if base is None else base

    out = io.StringIO()
    for line in io.StringIO(modelfile):
      command, _, args = line.partition(' ')
      if command.upper() not in ['FROM', 'ADAPTER']:
        print(line, end='', file=out)
        continue

      path = Path(args.strip()).expanduser()
      path = path if path.is_absolute() else base / path
      if path.exists():
        args = f'@{self._create_blob(path)}\n'
      print(command, args, end='', file=out)

    return out.getvalue()

  def _create_blob(self, path: Union[str, Path]) -> str:
    sha256sum = sha256()
    with open(path, 'rb') as r:
      while True:
        chunk = r.read(32 * 1024)
        if not chunk:
          break
        sha256sum.update(chunk)

    digest = f'sha256:{sha256sum.hexdigest()}'

    try:
      self._request('HEAD', f'/api/blobs/{digest}')
    except ResponseError as e:
      if e.status_code != 404:
        raise

      with open(path, 'rb') as r:
        self._request('POST', f'/api/blobs/{digest}', content=r)

    return digest

  def delete(self, model: str) -> Mapping[str, Any]:
    response = self._request('DELETE', '/api/delete', json={'name': model})
    return {'status': 'success' if response.status_code == 200 else 'error'}

  def list(self) -> Mapping[str, Any]:
    return self._request('GET', '/api/tags').json()

  def copy(self, source: str, destination: str) -> Mapping[str, Any]:
    response = self._request('POST', '/api/copy', json={'source': source, 'destination': destination})
    return {'status': 'success' if response.status_code == 200 else 'error'}

  def show(self, model: str) -> Mapping[str, Any]:
    return self._request('POST', '/api/show', json={'name': model}).json()


class AsyncClient(BaseClient):
  def __init__(self, host: Optional[str] = None, **kwargs) -> None:
    super().__init__(httpx.AsyncClient, host, **kwargs)

  async def _request(self, method: str, url: str, **kwargs) -> httpx.Response:
    response = await self._client.request(method, url, **kwargs)

    try:
      response.raise_for_status()
    except httpx.HTTPStatusError as e:
      raise ResponseError(e.response.text, e.response.status_code) from None

    return response

  async def _stream(self, method: str, url: str, **kwargs) -> AsyncIterator[Mapping[str, Any]]:
    async def inner():
      async with self._client.stream(method, url, **kwargs) as r:
        try:
          r.raise_for_status()
        except httpx.HTTPStatusError as e:
          e.response.read()
          raise ResponseError(e.response.text, e.response.status_code) from None

        async for line in r.aiter_lines():
          partial = json.loads(line)
          if e := partial.get('error'):
            raise ResponseError(e)
          yield partial

    return inner()

  async def _request_stream(
    self,
    *args,
    stream: bool = False,
    **kwargs,
  ) -> Union[Mapping[str, Any], AsyncIterator[Mapping[str, Any]]]:
    if stream:
      return await self._stream(*args, **kwargs)

    response = await self._request(*args, **kwargs)
    return response.json()

  async def generate(
    self,
    model: str = '',
    prompt: str = '',
    system: str = '',
    template: str = '',
    context: Optional[Sequence[int]] = None,
    stream: bool = False,
    raw: bool = False,
    format: Literal['', 'json'] = '',
    images: Optional[Sequence[AnyStr]] = None,
    options: Optional[Options] = None,
    keep_alive: Optional[Union[float, str]] = None,
  ) -> Union[Mapping[str, Any], AsyncIterator[Mapping[str, Any]]]:
    """
    Create a response using the requested model.

    Raises `RequestError` if a model is not provided.

    Raises `ResponseError` if the request could not be fulfilled.

    Returns `GenerateResponse` if `stream` is `False`, otherwise returns an asynchronous `GenerateResponse` generator.
    """
    if not model:
      raise RequestError('must provide a model')

    return await self._request_stream(
      'POST',
      '/api/generate',
      json={
        'model': model,
        'prompt': prompt,
        'system': system,
        'template': template,
        'context': context or [],
        'stream': stream,
        'raw': raw,
        'images': [_encode_image(image) for image in images or []],
        'format': format,
        'options': options or {},
        'keep_alive': keep_alive,
      },
      stream=stream,
    )

  async def chat(
    self,
    model: str = '',
    messages: Optional[Sequence[Message]] = None,
    stream: bool = False,
    format: Literal['', 'json'] = '',
    options: Optional[Options] = None,
    keep_alive: Optional[Union[float, str]] = None,
  ) -> Union[Mapping[str, Any], AsyncIterator[Mapping[str, Any]]]:
    """
    Create a chat response using the requested model.

    Raises `RequestError` if a model is not provided.

    Raises `ResponseError` if the request could not be fulfilled.

    Returns `ChatResponse` if `stream` is `False`, otherwise returns an asynchronous `ChatResponse` generator.
    """
    if not model:
      raise RequestError('must provide a model')

    for message in messages or []:
      if not isinstance(message, dict):
        raise TypeError('messages must be a list of strings')
      if not (role := message.get('role')) or role not in ['system', 'user', 'assistant']:
        raise RequestError('messages must contain a role and it must be one of "system", "user", or "assistant"')
      if not message.get('content'):
        raise RequestError('messages must contain content')
      if images := message.get('images'):
        message['images'] = [_encode_image(image) for image in images]

    return await self._request_stream(
      'POST',
      '/api/chat',
      json={
        'model': model,
        'messages': messages,
        'stream': stream,
        'format': format,
        'options': options or {},
        'keep_alive': keep_alive,
      },
      stream=stream,
    )

  async def embeddings(
    self,
    model: str = '',
    prompt: str = '',
    options: Optional[Options] = None,
    keep_alive: Optional[Union[float, str]] = None,
  ) -> Sequence[float]:
    response = await self._request(
      'POST',
      '/api/embeddings',
      json={
        'model': model,
        'prompt': prompt,
        'options': options or {},
        'keep_alive': keep_alive,
      },
    )

    return response.json()

  async def pull(
    self,
    model: str,
    insecure: bool = False,
    stream: bool = False,
  ) -> Union[Mapping[str, Any], AsyncIterator[Mapping[str, Any]]]:
    """
    Raises `ResponseError` if the request could not be fulfilled.

    Returns `ProgressResponse` if `stream` is `False`, otherwise returns a `ProgressResponse` generator.
    """
    return await self._request_stream(
      'POST',
      '/api/pull',
      json={
        'name': model,
        'insecure': insecure,
        'stream': stream,
      },
      stream=stream,
    )

  async def push(
    self,
    model: str,
    insecure: bool = False,
    stream: bool = False,
  ) -> Union[Mapping[str, Any], AsyncIterator[Mapping[str, Any]]]:
    """
    Raises `ResponseError` if the request could not be fulfilled.

    Returns `ProgressResponse` if `stream` is `False`, otherwise returns a `ProgressResponse` generator.
    """
    return await self._request_stream(
      'POST',
      '/api/push',
      json={
        'name': model,
        'insecure': insecure,
        'stream': stream,
      },
      stream=stream,
    )

  async def create(
    self,
    model: str,
    path: Optional[Union[str, PathLike]] = None,
    modelfile: Optional[str] = None,
    stream: bool = False,
  ) -> Union[Mapping[str, Any], AsyncIterator[Mapping[str, Any]]]:
    """
    Raises `ResponseError` if the request could not be fulfilled.

    Returns `ProgressResponse` if `stream` is `False`, otherwise returns a `ProgressResponse` generator.
    """
    if (realpath := _as_path(path)) and realpath.exists():
      modelfile = await self._parse_modelfile(realpath.read_text(), base=realpath.parent)
    elif modelfile:
      modelfile = await self._parse_modelfile(modelfile)
    else:
      raise RequestError('must provide either path or modelfile')

    return await self._request_stream(
      'POST',
      '/api/create',
      json={
        'name': model,
        'modelfile': modelfile,
        'stream': stream,
      },
      stream=stream,
    )

  async def _parse_modelfile(self, modelfile: str, base: Optional[Path] = None) -> str:
    base = Path.cwd() if base is None else base

    out = io.StringIO()
    for line in io.StringIO(modelfile):
      command, _, args = line.partition(' ')
      if command.upper() not in ['FROM', 'ADAPTER']:
        print(line, end='', file=out)
        continue

      path = Path(args.strip()).expanduser()
      path = path if path.is_absolute() else base / path
      if path.exists():
        args = f'@{await self._create_blob(path)}\n'
      print(command, args, end='', file=out)

    return out.getvalue()

  async def _create_blob(self, path: Union[str, Path]) -> str:
    sha256sum = sha256()
    with open(path, 'rb') as r:
      while True:
        chunk = r.read(32 * 1024)
        if not chunk:
          break
        sha256sum.update(chunk)

    digest = f'sha256:{sha256sum.hexdigest()}'

    try:
      await self._request('HEAD', f'/api/blobs/{digest}')
    except ResponseError as e:
      if e.status_code != 404:
        raise

      async def upload_bytes():
        with open(path, 'rb') as r:
          while True:
            chunk = r.read(32 * 1024)
            if not chunk:
              break
            yield chunk

      await self._request('POST', f'/api/blobs/{digest}', content=upload_bytes())

    return digest

  async def delete(self, model: str) -> Mapping[str, Any]:
    response = await self._request('DELETE', '/api/delete', json={'name': model})
    return {'status': 'success' if response.status_code == 200 else 'error'}

  async def list(self) -> Mapping[str, Any]:
    response = await self._request('GET', '/api/tags')
    return response.json()

  async def copy(self, source: str, destination: str) -> Mapping[str, Any]:
    response = await self._request('POST', '/api/copy', json={'source': source, 'destination': destination})
    return {'status': 'success' if response.status_code == 200 else 'error'}

  async def show(self, model: str) -> Mapping[str, Any]:
    response = await self._request('POST', '/api/show', json={'name': model})
    return response.json()


def _encode_image(image) -> str:
  """
  >>> _encode_image(b'ollama')
  'b2xsYW1h'
  >>> _encode_image(io.BytesIO(b'ollama'))
  'b2xsYW1h'
  >>> _encode_image('LICENSE')
  'TUlUIExpY2Vuc2UKCkNvcHlyaWdodCAoYykgT2xsYW1hCgpQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5Cm9mIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlICJTb2Z0d2FyZSIpLCB0byBkZWFsCmluIHRoZSBTb2Z0d2FyZSB3aXRob3V0IHJlc3RyaWN0aW9uLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uIHRoZSByaWdodHMKdG8gdXNlLCBjb3B5LCBtb2RpZnksIG1lcmdlLCBwdWJsaXNoLCBkaXN0cmlidXRlLCBzdWJsaWNlbnNlLCBhbmQvb3Igc2VsbApjb3BpZXMgb2YgdGhlIFNvZnR3YXJlLCBhbmQgdG8gcGVybWl0IHBlcnNvbnMgdG8gd2hvbSB0aGUgU29mdHdhcmUgaXMKZnVybmlzaGVkIHRvIGRvIHNvLCBzdWJqZWN0IHRvIHRoZSBmb2xsb3dpbmcgY29uZGl0aW9uczoKClRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIHNoYWxsIGJlIGluY2x1ZGVkIGluIGFsbApjb3BpZXMgb3Igc3Vic3RhbnRpYWwgcG9ydGlvbnMgb2YgdGhlIFNvZnR3YXJlLgoKVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEICJBUyBJUyIsIFdJVEhPVVQgV0FSUkFOVFkgT0YgQU5ZIEtJTkQsIEVYUFJFU1MgT1IKSU1QTElFRCwgSU5DTFVESU5HIEJVVCBOT1QgTElNSVRFRCBUTyBUSEUgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFksCkZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFORCBOT05JTkZSSU5HRU1FTlQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRQpBVVRIT1JTIE9SIENPUFlSSUdIVCBIT0xERVJTIEJFIExJQUJMRSBGT1IgQU5ZIENMQUlNLCBEQU1BR0VTIE9SIE9USEVSCkxJQUJJTElUWSwgV0hFVEhFUiBJTiBBTiBBQ1RJT04gT0YgQ09OVFJBQ1QsIFRPUlQgT1IgT1RIRVJXSVNFLCBBUklTSU5HIEZST00sCk9VVCBPRiBPUiBJTiBDT05ORUNUSU9OIFdJVEggVEhFIFNPRlRXQVJFIE9SIFRIRSBVU0UgT1IgT1RIRVIgREVBTElOR1MgSU4gVEhFClNPRlRXQVJFLgo='
  >>> _encode_image(Path('LICENSE'))
  'TUlUIExpY2Vuc2UKCkNvcHlyaWdodCAoYykgT2xsYW1hCgpQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5Cm9mIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlICJTb2Z0d2FyZSIpLCB0byBkZWFsCmluIHRoZSBTb2Z0d2FyZSB3aXRob3V0IHJlc3RyaWN0aW9uLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uIHRoZSByaWdodHMKdG8gdXNlLCBjb3B5LCBtb2RpZnksIG1lcmdlLCBwdWJsaXNoLCBkaXN0cmlidXRlLCBzdWJsaWNlbnNlLCBhbmQvb3Igc2VsbApjb3BpZXMgb2YgdGhlIFNvZnR3YXJlLCBhbmQgdG8gcGVybWl0IHBlcnNvbnMgdG8gd2hvbSB0aGUgU29mdHdhcmUgaXMKZnVybmlzaGVkIHRvIGRvIHNvLCBzdWJqZWN0IHRvIHRoZSBmb2xsb3dpbmcgY29uZGl0aW9uczoKClRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIHNoYWxsIGJlIGluY2x1ZGVkIGluIGFsbApjb3BpZXMgb3Igc3Vic3RhbnRpYWwgcG9ydGlvbnMgb2YgdGhlIFNvZnR3YXJlLgoKVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEICJBUyBJUyIsIFdJVEhPVVQgV0FSUkFOVFkgT0YgQU5ZIEtJTkQsIEVYUFJFU1MgT1IKSU1QTElFRCwgSU5DTFVESU5HIEJVVCBOT1QgTElNSVRFRCBUTyBUSEUgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFksCkZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFORCBOT05JTkZSSU5HRU1FTlQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRQpBVVRIT1JTIE9SIENPUFlSSUdIVCBIT0xERVJTIEJFIExJQUJMRSBGT1IgQU5ZIENMQUlNLCBEQU1BR0VTIE9SIE9USEVSCkxJQUJJTElUWSwgV0hFVEhFUiBJTiBBTiBBQ1RJT04gT0YgQ09OVFJBQ1QsIFRPUlQgT1IgT1RIRVJXSVNFLCBBUklTSU5HIEZST00sCk9VVCBPRiBPUiBJTiBDT05ORUNUSU9OIFdJVEggVEhFIFNPRlRXQVJFIE9SIFRIRSBVU0UgT1IgT1RIRVIgREVBTElOR1MgSU4gVEhFClNPRlRXQVJFLgo='
  >>> _encode_image('YWJj')
  'YWJj'
  >>> _encode_image(b'YWJj')
  'YWJj'
  """

  if p := _as_path(image):
    return b64encode(p.read_bytes()).decode('utf-8')

  try:
    b64decode(image, validate=True)
    return image if isinstance(image, str) else image.decode('utf-8')
  except (binascii.Error, TypeError):
    ...

  if b := _as_bytesio(image):
    return b64encode(b.read()).decode('utf-8')

  raise RequestError('image must be bytes, path-like object, or file-like object')


def _as_path(s: Optional[Union[str, PathLike]]) -> Union[Path, None]:
  if isinstance(s, str) or isinstance(s, Path):
    try:
      if (p := Path(s)).exists():
        return p
    except Exception:
      ...
  return None


def _as_bytesio(s: Any) -> Union[io.BytesIO, None]:
  if isinstance(s, io.BytesIO):
    return s
  elif isinstance(s, bytes):
    return io.BytesIO(s)
  return None


def _parse_host(host: Optional[str]) -> str:
  """
  >>> _parse_host(None)
  'http://127.0.0.1:11434'
  >>> _parse_host('')
  'http://127.0.0.1:11434'
  >>> _parse_host('1.2.3.4')
  'http://1.2.3.4:11434'
  >>> _parse_host(':56789')
  'http://127.0.0.1:56789'
  >>> _parse_host('1.2.3.4:56789')
  'http://1.2.3.4:56789'
  >>> _parse_host('http://1.2.3.4')
  'http://1.2.3.4:80'
  >>> _parse_host('https://1.2.3.4')
  'https://1.2.3.4:443'
  >>> _parse_host('https://1.2.3.4:56789')
  'https://1.2.3.4:56789'
  >>> _parse_host('example.com')
  'http://example.com:11434'
  >>> _parse_host('example.com:56789')
  'http://example.com:56789'
  >>> _parse_host('http://example.com')
  'http://example.com:80'
  >>> _parse_host('https://example.com')
  'https://example.com:443'
  >>> _parse_host('https://example.com:56789')
  'https://example.com:56789'
  >>> _parse_host('example.com/')
  'http://example.com:11434'
  >>> _parse_host('example.com:56789/')
  'http://example.com:56789'
  """

  host, port = host or '', 11434
  scheme, _, hostport = host.partition('://')
  if not hostport:
    scheme, hostport = 'http', host
  elif scheme == 'http':
    port = 80
  elif scheme == 'https':
    port = 443

  split = urllib.parse.urlsplit('://'.join([scheme, hostport]))
  host = split.hostname or '127.0.0.1'
  port = split.port or port

  return f'{scheme}://{host}:{port}'


_client = Client()

generate = _client.generate
chat = _client.chat
embeddings = _client.embeddings
pull = _client.pull
push = _client.push
create = _client.create
delete = _client.delete
list = _client.list
copy = _client.copy
show = _client.show
