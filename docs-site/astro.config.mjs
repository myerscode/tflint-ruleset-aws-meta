// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
  site: 'https://myerscode.github.io',
  base: '/tflint-ruleset-aws-meta/',

  integrations: [
    starlight({
      title: 'TFLint AWS Meta Ruleset',
      description: 'Documentation for the tflint-ruleset-aws-meta project',
      sidebar: [
        {
          label: 'Overview',
          items: [
            { label: 'Introduction', link: '/' },
            { label: 'Installation', link: '/installation' },
            { label: 'Configuration', link: '/configuration' },
            // { label: 'Contributing', link: '/contributing' },
          ],
        },
        {
          label: 'Rules',
          autogenerate: { directory: 'rules' },
        },
      ],
    }),
  ],
});
