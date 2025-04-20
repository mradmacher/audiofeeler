require "spec"
require "webdrivers"
require "selenium"

macro try!(expr)
  begin
    {{expr}}
  rescue e
    fail(e.message)
  end
end

class MyWebdriver
  getter driver : Selenium::Driver

  def initialize
    webdriver_path = Webdrivers::Chromedriver.install
    service = Selenium::Service.chrome(driver_path: webdriver_path)
    @driver = Selenium::Driver.for(:chrome, service: service)

    @capabilities = Selenium::Chrome::Capabilities.new
    @capabilities.chrome_options.args = ["no-sandbox", "headless", "disable-gpu", "user-data-dir=#{__DIR__}../data/selenium"]

  end

  def create_session
    @driver.create_session(@capabilities)
  end

  def stop
    @driver.stop
  end
end

DRIVER = MyWebdriver.new

Spec.after_suite do
  DRIVER.stop
end

def with_session(&)
  session = DRIVER.create_session
  yield(session)
ensure
  session.delete unless session.nil?
end
