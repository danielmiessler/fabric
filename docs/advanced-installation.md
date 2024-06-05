# Advanced Installation

## Table of Contents
- [Using `fabric` with Inference Server](#using-fabric-with-inference-server)


## Using `fabric` with Inference Server
If you want to use it with OpenAI API compatible inference servers, such as [FastChat](https://github.com/lm-sys/FastChat), [Helmholtz Blablador](http://helmholtz-blablador.fz-juelich.de), [LM Studio](https://lmstudio.ai) and others, simply export the following environment variables:

- `export OPENAI_BASE_URL=https://YOUR-SERVER:8000/v1/`
- `export DEFAULT_MODEL="YOUR_MODEL"`

And if your server needs authentication tokens, like Blablador does, you export the token the same way you would with OpenAI:
  
- `export OPENAI_API_KEY="YOUR TOKEN"`

