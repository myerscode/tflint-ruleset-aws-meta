# Documentation Site

This directory contains the documentation site for the tflint-ruleset-aws-meta project, built with [Astro](https://astro.build/) and [Starlight](https://starlight.astro.build/).

## Development

Install dependencies:

```bash
npm install
```

Start the development server:

```bash
npm run dev
```

The site will be available at `http://localhost:4321/tflint-ruleset-aws-meta/`

## Building

Build the static site:

```bash
npm run build
```

Preview the built site:

```bash
npm run preview
```

## Deployment

The site is automatically deployed to GitHub Pages when changes are pushed to the `main` branch via the `.github/workflows/deploy-pages.yml` workflow.

The site is published at: https://myerscode.github.io/tflint-ruleset-aws-meta/

## Content Structure

- `src/content/docs/index.md` - Homepage/Introduction
- `src/content/docs/installation.md` - Installation and configuration guide
- `src/content/docs/contributing.md` - Development and contribution guide
- `src/content/docs/rules/` - Individual rule documentation pages
- `astro.config.mjs` - Site configuration including navigation

## Adding New Rules

When adding a new rule to the ruleset:

1. Create a new markdown file in `src/content/docs/rules/` with the rule name (e.g., `aws_new_rule.md`)
2. Add frontmatter with:
   - `title`: User-friendly readable title
   - `description`: Brief description of what the rule does
   - `ruleName`: The actual rule name (e.g., `aws_new_rule`)
3. Include the rule name prominently in the document body: `**Rule:** \`aws_new_rule\``
4. Include example violations and recommended fixes
5. Update `src/content/docs/rules/index.md` to include the new rule with its user-friendly title
6. Update the rules table in `src/content/docs/index.md`
7. Update the rules table in the main `README.md`
