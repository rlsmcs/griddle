project:
  name: griddle
  description: formula 1 guess the driver minigame
  status: deployed
  url: https://griddle-app.vercel.app/

tech_stack:
  frontend:
    - next.js
    - tailwind css
    - typescript
  backend:
    - golang
    - rest api

structure:
  root:
    app:
      - globals.css
      - layout.tsx
      - page.tsx
    components:
      - ColumnHeaders.tsx
      - GameOverModal.tsx
      - GuessRow.tsx
      - HintTile.tsx
      - Legend.tsx
      - SearchInput.tsx
    data:
      - drivers.json
      - drivers_sample.json
    public:
      data:
        - drivers.json
    lib:
      - game.ts
      - types.ts
    scripts: []
    config:
      - next.config.js
      - postcss.config.js
      - tailwind.config.js
      - tsconfig.json
      - next-env.d.ts
    other:
      - package.json
      - package-lock.json
      - README.md
      - LICENSE
      - .gitignore
      - architecture.yaml

deployment:
  platform: vercel
  type: frontend hosting
  backend: external go service (used for dataset updates)

notes:
  - static driver dataset generated using a golang script
  - dataset updated on a race-by-race basis with latest statistics again using the same golang script
  - supports adding new drivers (currently includes data from 2018 onwards)
  - game compares attributes such as nationality, team, age, debut, and wins
  - six guess limit per session
  - simple color and directional hint system
  - color logic will be improved (e.g. better relative comparisons following the purple/green/yellow sector time colours in f1)
  - ui is minimal, not to scale, and will be improved later
