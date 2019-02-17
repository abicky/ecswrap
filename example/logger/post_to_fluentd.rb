require 'logger'
require 'fluent-logger'

$stdout.sync = true
logger = Logger.new($stdout)

stop = false
trap :SIGTERM do
  Thread.new do
    logger.info 'Stopping...'
    sleep 5
    stop = true
  end
end

log = Fluent::Logger::FluentLogger.new(nil, host: 'fluentd')
cnt = 0
until stop
  cnt += 1
  logger.info "Post an event (cnt: #{cnt})"
  unless log.post('', cnt: cnt)
    logger.error "Failed to post an event (cnt: #{cnt})"
  end
  sleep 1
end

logger.info 'Stopped'
