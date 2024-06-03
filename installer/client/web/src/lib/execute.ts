import { exec } from "child_process"

export type ExecuteOutput = { ok: false, error: string[], command: string } | { ok: true, data: string[], command: string }

export async function execute(command: string) {
  const output = await new Promise<ExecuteOutput>((resolve, reject) => {
    exec(command, function (error, stdout, stderr) {
      if (error) {
        reject({ ok: false, error: stderr.toString().split("\n"), command })
      }
      resolve({ ok: true, data: stdout.toString().split("\n"), command })
    });
  })
  console.log({ output, type: typeof output })
  return output
};
