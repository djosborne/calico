---
title: Release notes
canonical_url: 'https://docs.projectcalico.org/v3.7/release-notes/index'
---

The following table shows component versioning for {{site.prodname}}  **{{ page.version }}**.

Use the version selector at the top-right of this page to view a different release.

{% for release in site.data.versions[page.version] %}
## {{ release.title }}
{% unless release.title == "master" %}
[Release archive](https://github.com/projectcalico/calico/releases/download/{{ release.title }}/release-{{ release.title }}.tgz) with Kubernetes manifests, Docker images and binaries.
{% endunless %}

{% if release.note %}
{{ release.note }}
{% else %}
{% include {{page.version}}/release-notes/{{release.title}}-release-notes.md %}
{% endif %}

| Component              | Version |
|------------------------|---------|
{% for component in release.components %}
{%- assign name = component[0] -%}
{%- assign version = component[1].version -%}

| {{ name }} | [{{ version }}]({% include component_url component=name version=version %}) |
{% endfor %}
{% endfor %}

