require "./../browser_spec_helper"

describe "selenium" do
  it "works" do
    with_session do |session|
      session.navigate_to("http://localhost:3000/accounts")
      #element = session.find_element(:link_text, "Etnorakieta")
      element = try!(session.find_element(:link_text, "Etnorakieta"))
      element.click

      session.title.should eq("Audiofeeler")
    end
  end

  it "also works" do
    with_session do |session|
      session.navigate_to("http://localhost:3000/accounts")
      element = session.find_element(:link_text, "Etnorozrabiaka")
      element.click

      session.find_element(:tag_name, "h1").text.should eq "Etnorozrabiaka"
    end
  end
end
