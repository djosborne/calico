require "jekyll"
require "tempfile"

module Jekyll
  class RenderHelmTagBlock < Liquid::Block
    def render(context)
      text = super

      # Because helm hasn't merged stdin support, write our values to a tempfile fml
      t = Tempfile.new("jhelm")
      t.write(text)
      t.close

      # TODO: load the version from page.version
      version = "master"
      out = `helm template _includes/#{version}/charts/calico -f #{t.path}`
      
      t.unlink
      return out
    end
  end
end

Liquid::Template.register_tag('helm', Jekyll::RenderHelmTagBlock)
