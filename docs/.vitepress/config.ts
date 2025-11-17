import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Studio MCP',
  description: 'Turn any CLI command into an MCP tool',

  themeConfig: {
    logo: '/logo.svg',

    nav: [
      { text: 'Guide', link: '/guide/getting-started' },
      { text: 'Reference', link: '/reference/template-syntax' },
      { text: 'Examples', link: '/examples/basic' },
      { text: 'GitHub', link: 'https://github.com/studio-mcp/studio' }
    ],

    sidebar: {
      '/guide/': [
        {
          text: 'Introduction',
          items: [
            { text: 'What is Studio MCP?', link: '/guide/what-is-studio' },
            { text: 'Getting Started', link: '/guide/getting-started' },
            { text: 'Installation', link: '/guide/installation' }
          ]
        },
        {
          text: 'Configuration',
          items: [
            { text: 'Claude Desktop', link: '/guide/config-claude' },
            { text: 'Cursor', link: '/guide/config-cursor' },
            { text: 'VSCode', link: '/guide/config-vscode' }
          ]
        },
        {
          text: 'Advanced',
          items: [
            { text: 'Debugging', link: '/guide/debugging' },
            { text: 'Development', link: '/guide/development' }
          ]
        }
      ],
      '/reference/': [
        {
          text: 'Reference',
          items: [
            { text: 'Template Syntax', link: '/reference/template-syntax' },
            { text: 'Command Line Options', link: '/reference/cli-options' },
            { text: 'Architecture', link: '/reference/architecture' }
          ]
        }
      ],
      '/examples/': [
        {
          text: 'Examples',
          items: [
            { text: 'Basic Examples', link: '/examples/basic' },
            { text: 'Common Tools', link: '/examples/common-tools' },
            { text: 'Advanced Use Cases', link: '/examples/advanced' }
          ]
        }
      ]
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/studio-mcp/studio' }
    ],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2025 Martin Emde'
    },

    search: {
      provider: 'local'
    }
  }
})
