autoload :ChildProcess, 'childprocess'
autoload :Tempfile, 'tempfile'

module ChildProcessHelper
  module_function def call(cmd, env: nil, cwd: nil)
    process = ChildProcess.new(*cmd)
    process.io.inherit!
    if cwd
      process.cwd = cwd
    end
    if env
      env.each do |k, v|
        process.environment[k.to_s] = v
      end
    end
    process.start
    process.wait
    process
  end

  module_function def check_call(cmd, env: nil, cwd: nil)
    process = call(cmd, env: env, cwd: cwd)
    unless process.exit_code == 0
      raise "Failed to execute: #{cmd}"
    end
  end

  module_function def check_output(cmd, env: nil, cwd: nil)
    process = ChildProcess.new(*cmd)
    process.io.inherit!
    process.io.stdout = Tempfile.new("child-output")
    if cwd
      process.cwd = cwd
    end
    if env
      env.each do |k, v|
        process.environment[k.to_s] = v
      end
    end
    process.start
    process.wait
    unless process.exit_code == 0
      raise "Failed to execute: #{cmd}"
    end
    process.io.stdout.rewind
    process.io.stdout.read
  end
end
