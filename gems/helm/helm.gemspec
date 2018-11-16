Gem::Specification.new do |s|
  s.name = %q{helm}
  s.version = "0.0.1"
  s.date = %q{2018-11-29}
  s.summary = %q{render helm templates in jekyll}
  s.author = %q{dan@projectcalico.org}
  s.files = [
    "Gemfile",
    "Rakefile",
    "lib/helm.rb"
  ]
  s.require_paths = ["lib"]
  s.add_dependency 'jekyll', "3.7.4"
end
