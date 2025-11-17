# Studio MCP Documentation

This directory contains the VitePress documentation site for Studio MCP.

## Getting Started

### Install Dependencies

```bash
cd docs
npm install
```

### Development Server

Run the local development server:

```bash
npm run dev
```

The site will be available at `http://localhost:5173`

### Build for Production

Build the static site:

```bash
npm run build
```

Output will be in `.vitepress/dist/`

### Preview Production Build

Preview the production build locally:

```bash
npm run preview
```

## Documentation Structure

```
docs/
├── .vitepress/
│   └── config.ts           # VitePress configuration
├── guide/                   # User guides
│   ├── what-is-studio.md
│   ├── getting-started.md
│   ├── installation.md
│   ├── config-claude.md
│   ├── config-cursor.md
│   ├── config-vscode.md
│   ├── debugging.md
│   └── development.md
├── reference/              # API and technical reference
│   ├── template-syntax.md
│   ├── cli-options.md
│   └── architecture.md
├── examples/               # Examples and tutorials
│   ├── basic.md
│   ├── common-tools.md
│   └── advanced.md
└── index.md               # Home page
```

## Adding Content

### New Guide Page

1. Create a new `.md` file in the appropriate directory
2. Add frontmatter if needed
3. Update `.vitepress/config.ts` sidebar configuration

### New Example

1. Add example to `examples/` directory
2. Update sidebar in `.vitepress/config.ts`
3. Link from related pages

## Writing Tips

- Use clear, concise language
- Include code examples for every feature
- Add links to related sections
- Test all code examples
- Use proper markdown formatting

## VitePress Features

- **Markdown Extensions**: Enhanced markdown with custom containers
- **Code Blocks**: Syntax highlighting and line numbers
- **Search**: Built-in search functionality
- **Mobile Responsive**: Works on all devices
- **Fast**: Static site generation for performance

## Deployment

The documentation can be deployed to:

- GitHub Pages
- Netlify
- Vercel
- Any static hosting service

Just build the site and deploy the `.vitepress/dist` directory.

## Contributing

When contributing to the documentation:

1. Follow the existing structure
2. Maintain consistent tone and style
3. Test all code examples
4. Check for broken links
5. Preview changes locally before committing
