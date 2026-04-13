```yaml
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
  backend: external go service

notes:
  - static driver dataset generated using a golang script
  - dataset updated on a race-by-race basis using same golang script
  - supports adding new drivers (2018 onwards)
  - attribute comparison: nationality, team, age, debut, wins
  - six guess limit
  - simple color + directional hint system (later will update to correlate with the f1 sector times (purple,yellow,green,blue,black etc.))
  - ui is minimal and will be improved
