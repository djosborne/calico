require "jekyll"
require "tempfile"
require "open3"

require_relative "./lib"

# This plugin enables jekyll to render helm charts.
# Traditionally, Jekyll will render files which make use of the Liquid templating language.
# This plugin adds a new 'tag' that when specified will pass the input to the Helm binary.
# example use:
#
# {% helm %}
# datastore: kubernetes
# networking: calico
# {% endhelm %}
module Jekyll
  class RenderHelmTagBlock < Liquid::Block
    def initialize(tag_name, extra_args, liquid_options)
      super

      @chart = "calico"
      if extra_args.start_with?("tigera-operator")
        @chart = "tigera-operator"
        extra_args.slice! "tigera-operator"
      end

      # substitute --execute with --show-only for helm v3 compatibility.
      extra_args.gsub!(/--execute (\S*)/) do |f|
        # calico CRDs stay in the templates/crds directory
        if $1.start_with? "templates/crds/calico" then return f end
        # operator CRDs have moved to root
        if $1.start_with? "templates/crds/" then return f.gsub('--execute templates/crds/', '--show-only ') end
        # all other requests need to use --show-only instead of --execute for helm v3
        return f.gsub('--execute', '--show-only')
      end

      @extra_args = extra_args
    end
    def render(context)
      text = super

      # Because helm hasn't merged stdin support, write the passed-in values.yaml
      # to a tempfile on disk.
      t = Tempfile.new("jhelm")
      t.write(text)
      t.close

      imageRegistry = context.registers[:page]["registry"]
      imageNames = context.registers[:site].config["imageNames"]
      versions = context.registers[:site].data["versions"]

      vs = parse_versions(versions)

      versionsYml = gen_values(vs, imageNames, imageRegistry, @chart)

      tv = Tempfile.new("temp_versions.yml")
      tv.write(versionsYml)
      tv.close

      # Execute helm.
      # Set the default etcd endpoint placeholder for rendering in the docs.
      cmd = """helm template --include-crds _includes/charts/#{@chart} \
        -f #{tv.path} \
        -f #{t.path} \
        --set etcd.endpoints='http://<ETCD_IP>:<ETCD_PORT>'"""

      cmd += " " + @extra_args.to_s

      out, stderr, status = Open3.capture3(cmd)
      if status != 0
        raise "failed to execute helm for '#{context.registers[:page]["path"]}': #{stderr}"
      end

      t.unlink
      tv.unlink
      return out
    end
  end
end

Liquid::Template.register_tag('helm', Jekyll::RenderHelmTagBlock)
