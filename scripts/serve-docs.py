#!/usr/bin/env python3
"""
Local Terraform Registry-style documentation server.
Serves provider docs with a layout similar to registry.terraform.io.

Usage: python3 scripts/serve-docs.py [port]
Default port: 8081
"""

import http.server
import json
import os
import re
import sys
from pathlib import Path
from urllib.parse import unquote

PORT = int(sys.argv[1]) if len(sys.argv) > 1 else 8081
DOCS_DIR = Path(__file__).resolve().parent.parent / "docs"
PROVIDER_NAME = "gcore"

def list_docs(subdir):
    d = DOCS_DIR / subdir
    if not d.exists():
        return []
    items = []
    for f in sorted(d.iterdir()):
        if f.suffix == ".md":
            name = f.stem
            title = "gcore_" + name
            # Read first description line
            desc = ""
            with open(f) as fh:
                for line in fh:
                    if line.startswith("description:"):
                        desc = line.split(":", 1)[1].strip().strip('"').strip("|").strip()
                        break
            items.append({"name": name, "title": title, "desc": desc})
    return items

def read_doc(path):
    full = DOCS_DIR / path
    if not full.exists():
        return None
    with open(full) as f:
        return f.read()

def strip_frontmatter(md):
    if md.startswith("---"):
        end = md.find("---", 3)
        if end != -1:
            return md[end + 3:].strip()
    return md

RESOURCES = list_docs("resources")
DATA_SOURCES = list_docs("data-sources")

def build_sidebar():
    html = '<div class="sidebar-section"><h3>Resources</h3><ul>'
    for r in RESOURCES:
        html += '<li><a href="/resources/{name}" class="nav-link" data-path="resources/{name}.md">{title}</a></li>'.format(**r)
    html += '</ul></div><div class="sidebar-section"><h3>Data Sources</h3><ul>'
    for d in DATA_SOURCES:
        html += '<li><a href="/data-sources/{name}" class="nav-link" data-path="data-sources/{name}.md">{title}</a></li>'.format(**d)
    html += '</ul></div>'
    return html

SIDEBAR_HTML = build_sidebar()

PAGE_TEMPLATE = """<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{title} - Gcore Provider - Terraform Registry</title>
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github.min.css">
<style>
:root {{
  --sidebar-w: 280px;
  --purple: #7b42bc;
  --purple-dark: #5c2d91;
  --border: #e2e8f0;
  --bg: #f7f8fa;
  --text: #1a202c;
  --text-secondary: #64748b;
  --link: #7b42bc;
}}
* {{ margin: 0; padding: 0; box-sizing: border-box; }}
body {{ font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif; color: var(--text); background: var(--bg); }}

/* Top bar */
.topbar {{ background: #000; color: #fff; padding: 0 24px; height: 52px; display: flex; align-items: center; gap: 16px; position: fixed; top: 0; left: 0; right: 0; z-index: 100; }}
.topbar .logo {{ font-weight: 700; font-size: 1rem; display: flex; align-items: center; gap: 8px; }}
.topbar .logo svg {{ width: 20px; height: 20px; }}
.topbar .breadcrumb {{ font-size: 0.85rem; color: #a0aec0; }}
.topbar .breadcrumb a {{ color: #a0aec0; text-decoration: none; }}
.topbar .breadcrumb a:hover {{ color: #fff; }}
.topbar .breadcrumb .sep {{ margin: 0 6px; }}

/* Provider header */
.provider-header {{ background: #fff; border-bottom: 1px solid var(--border); padding: 20px 24px 20px calc(var(--sidebar-w) + 24px); margin-top: 52px; }}
.provider-header h1 {{ font-size: 1.5rem; font-weight: 600; }}
.provider-header .meta {{ font-size: 0.85rem; color: var(--text-secondary); margin-top: 4px; }}

/* Layout */
.layout {{ display: flex; margin-top: 52px; }}
.sidebar {{ width: var(--sidebar-w); min-width: var(--sidebar-w); background: #fff; border-right: 1px solid var(--border); position: fixed; top: 52px; bottom: 0; overflow-y: auto; padding: 16px 0; z-index: 50; }}
.main {{ margin-left: var(--sidebar-w); flex: 1; min-width: 0; padding: 32px 48px 64px; max-width: 900px; }}

/* Sidebar */
.sidebar-section h3 {{ font-size: 0.7rem; text-transform: uppercase; letter-spacing: 0.08em; color: var(--text-secondary); padding: 12px 20px 6px; font-weight: 600; }}
.sidebar ul {{ list-style: none; }}
.sidebar li {{ }}
.sidebar .nav-link {{ display: block; padding: 4px 20px; font-size: 0.85rem; color: var(--text); text-decoration: none; border-left: 3px solid transparent; }}
.sidebar .nav-link:hover {{ background: #f1f5f9; color: var(--link); }}
.sidebar .nav-link.active {{ border-left-color: var(--purple); color: var(--purple); font-weight: 600; background: #faf5ff; }}

/* Search */
.sidebar-search {{ padding: 8px 16px 8px; }}
.sidebar-search input {{ width: 100%; padding: 6px 10px; border: 1px solid var(--border); border-radius: 6px; font-size: 0.85rem; outline: none; }}
.sidebar-search input:focus {{ border-color: var(--purple); box-shadow: 0 0 0 2px rgba(123,66,188,0.15); }}

/* Content */
.content h1 {{ font-size: 1.75rem; font-weight: 700; margin-bottom: 8px; padding-bottom: 12px; border-bottom: 1px solid var(--border); }}
.content h2 {{ font-size: 1.25rem; font-weight: 600; margin: 28px 0 12px; padding-bottom: 8px; border-bottom: 1px solid var(--border); }}
.content h3 {{ font-size: 1.05rem; font-weight: 600; margin: 20px 0 8px; }}
.content p {{ margin: 8px 0; line-height: 1.7; }}
.content ul, .content ol {{ margin: 8px 0 8px 24px; }}
.content li {{ margin: 4px 0; line-height: 1.6; }}
.content a {{ color: var(--link); text-decoration: none; }}
.content a:hover {{ text-decoration: underline; }}
.content code {{ background: #f1f5f9; padding: 2px 6px; border-radius: 4px; font-size: 0.9em; font-family: 'SF Mono', 'Fira Code', monospace; }}
.content pre {{ background: #1e293b; color: #e2e8f0; padding: 16px 20px; border-radius: 8px; overflow-x: auto; margin: 12px 0; line-height: 1.5; }}
.content pre code {{ background: none; color: inherit; padding: 0; font-size: 0.85rem; }}
.content blockquote {{ border-left: 4px solid var(--purple); padding: 8px 16px; margin: 12px 0; background: #faf5ff; border-radius: 0 6px 6px 0; }}
.content table {{ width: 100%; border-collapse: collapse; margin: 12px 0; }}
.content th {{ background: #f1f5f9; text-align: left; padding: 8px 12px; font-size: 0.85rem; font-weight: 600; border: 1px solid var(--border); }}
.content td {{ padding: 8px 12px; border: 1px solid var(--border); font-size: 0.9rem; vertical-align: top; }}

/* Page type badge */
.page-type {{ display: inline-block; padding: 2px 10px; border-radius: 4px; font-size: 0.75rem; font-weight: 600; text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 8px; }}
.page-type.resource {{ background: #dbeafe; color: #1e40af; }}
.page-type.datasource {{ background: #dcfce7; color: #166534; }}

/* Nested schema styling */
.content h4 {{ font-size: 0.95rem; font-weight: 600; margin: 16px 0 6px; }}

/* Required/Optional labels */
.content li code:first-child {{ font-weight: 600; }}
</style>
</head>
<body>

<div class="topbar">
  <div class="logo">
    <svg viewBox="0 0 256 256" fill="currentColor"><polygon points="128,0 256,74 256,182 128,256 0,182 0,74"/></svg>
    Terraform
  </div>
  <div class="breadcrumb">
    <a href="/">Registry</a><span class="sep">/</span>
    <a href="/">Providers</a><span class="sep">/</span>
    <a href="/">gcore</a><span class="sep">/</span>
    <span>{breadcrumb}</span>
  </div>
</div>

<div class="layout">
  <div class="sidebar">
    <div class="sidebar-search">
      <input type="text" id="search" placeholder="Filter resources..." oninput="filterSidebar(this.value)">
    </div>
    <a href="/" class="nav-link" style="padding: 8px 20px; font-weight: 600; border-bottom: 1px solid var(--border); margin-bottom: 4px; display: block;">Overview</a>
    {sidebar}
  </div>
  <div class="main">
    <div class="content" id="content">
      {content}
    </div>
  </div>
</div>

<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/languages/hcl.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/languages/bash.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/languages/shell.min.js"></script>
<script>
hljs.highlightAll();

// Highlight active sidebar link
var path = window.location.pathname;
document.querySelectorAll('.nav-link').forEach(function(a) {{
  if (a.getAttribute('href') === path) a.classList.add('active');
}});

function filterSidebar(q) {{
  q = q.toLowerCase();
  document.querySelectorAll('.sidebar .nav-link').forEach(function(a) {{
    var li = a.parentElement;
    if (!li || li.tagName !== 'LI') return;
    li.style.display = a.textContent.toLowerCase().includes(q) ? '' : 'none';
  }});
}}
</script>
</body>
</html>"""

def md_to_html(md):
    """Simple markdown to HTML converter for Terraform docs."""
    md = strip_frontmatter(md)
    lines = md.split('\n')
    html_parts = []
    in_code = False
    code_lang = ''
    code_lines = []
    in_list = False
    in_table = False
    table_lines = []

    def flush_table():
        nonlocal table_lines, in_table
        if not table_lines:
            return ''
        rows = []
        for i, row in enumerate(table_lines):
            cells = [c.strip() for c in row.strip('|').split('|')]
            if i == 1 and all(set(c.strip()) <= set('-: ') for c in cells):
                continue
            tag = 'th' if i == 0 else 'td'
            row_html = ''.join('<{0}>{1}</{0}>'.format(tag, process_inline(c)) for c in cells)
            rows.append('<tr>' + row_html + '</tr>')
        table_lines = []
        in_table = False
        return '<table>' + ''.join(rows) + '</table>'

    def process_inline(text):
        # Code
        text = re.sub(r'`([^`]+)`', r'<code>\1</code>', text)
        # Bold
        text = re.sub(r'\*\*(.+?)\*\*', r'<strong>\1</strong>', text)
        # Italic
        text = re.sub(r'\*(.+?)\*', r'<em>\1</em>', text)
        # Links
        text = re.sub(r'\[([^\]]+)\]\(([^)]+)\)', r'<a href="\2">\1</a>', text)
        return text

    for line in lines:
        # Code blocks
        if line.startswith('```'):
            if in_code:
                lang_class = ' class="language-{}"'.format(code_lang) if code_lang else ''
                html_parts.append('<pre><code{}>{}</code></pre>'.format(
                    lang_class,
                    '\n'.join(code_lines).replace('&', '&amp;').replace('<', '&lt;').replace('>', '&gt;')
                ))
                code_lines = []
                in_code = False
            else:
                if in_table:
                    html_parts.append(flush_table())
                code_lang = line[3:].strip()
                if code_lang == 'terraform' or code_lang == 'tf':
                    code_lang = 'hcl'
                in_code = True
            continue
        if in_code:
            code_lines.append(line)
            continue

        # Table
        if '|' in line and line.strip().startswith('|'):
            in_table = True
            table_lines.append(line)
            continue
        elif in_table:
            html_parts.append(flush_table())

        # Headers
        m = re.match(r'^(#{1,6})\s+(.+)', line)
        if m:
            if in_list:
                html_parts.append('</ul>')
                in_list = False
            level = len(m.group(1))
            html_parts.append('<h{0}>{1}</h{0}>'.format(level, process_inline(m.group(2))))
            continue

        # List items
        if re.match(r'^[-*]\s', line.strip()):
            if not in_list:
                html_parts.append('<ul>')
                in_list = True
            html_parts.append('<li>{}</li>'.format(process_inline(line.strip()[2:])))
            continue
        elif in_list and line.strip() == '':
            html_parts.append('</ul>')
            in_list = False
            continue

        # Blockquote
        if line.startswith('>'):
            html_parts.append('<blockquote>{}</blockquote>'.format(process_inline(line[1:].strip())))
            continue

        # Horizontal rule
        if re.match(r'^---+\s*$', line):
            html_parts.append('<hr>')
            continue

        # Paragraph
        stripped = line.strip()
        if stripped:
            html_parts.append('<p>{}</p>'.format(process_inline(stripped)))

    if in_list:
        html_parts.append('</ul>')
    if in_table:
        html_parts.append(flush_table())

    return '\n'.join(html_parts)


class DocsHandler(http.server.BaseHTTPRequestHandler):
    def log_message(self, format, *args):
        # Quieter logging
        pass

    def do_GET(self):
        path = unquote(self.path)

        # Route: /
        if path == '/' or path == '/index':
            md = read_doc('index.md')
            if md:
                self.serve_page('Gcore Provider', 'Documentation', md_to_html(md))
            else:
                self.send_error(404)
            return

        # Route: /resources/<name> or /data-sources/<name>
        m = re.match(r'^/(resources|data-sources)/([a-z0-9_]+)$', path)
        if m:
            category, name = m.groups()
            md = read_doc('{}/{}.md'.format(category, name))
            if md:
                type_label = 'Resource' if category == 'resources' else 'Data Source'
                badge_class = 'resource' if category == 'resources' else 'datasource'
                badge = '<span class="page-type {}">{}</span>'.format(badge_class, type_label)
                title = 'gcore_{} ({})'.format(name, type_label)
                breadcrumb = '{} / gcore_{}'.format(type_label + 's', name)
                self.serve_page(title, breadcrumb, badge + md_to_html(md))
            else:
                self.send_error(404)
            return

        self.send_error(404)

    def serve_page(self, title, breadcrumb, content):
        html = PAGE_TEMPLATE.format(
            title=title,
            breadcrumb=breadcrumb,
            sidebar=SIDEBAR_HTML,
            content=content
        )
        self.send_response(200)
        self.send_header('Content-Type', 'text/html; charset=utf-8')
        self.end_headers()
        self.wfile.write(html.encode('utf-8'))

if __name__ == '__main__':
    server = http.server.HTTPServer(('127.0.0.1', PORT), DocsHandler)
    print('Serving Terraform docs at http://localhost:{}'.format(PORT))
    print('  Resources: {}, Data Sources: {}'.format(len(RESOURCES), len(DATA_SOURCES)))
    print('  Press Ctrl+C to stop')
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print('\nStopped.')
