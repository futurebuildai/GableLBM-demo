# GableLBM: Design System & UX Principles

To disrupt legacy ERPs, GableLBM must feel like a premium, high-performance tool—more like a Bloomberg Terminal or a Tesla interface than a traditional business app.

## 1. The Aesthetic: "Industrial Dark"
We prioritize high contrast, deep depth, and focused visibility.

### Color Palette
- **Primary:** `#00FFA3` (Electric Lumber Green) - Used for primary actions and "success" states.
- **Background:** `#0A0B10` (Deep Space) - The main canvas.
- **Surface:** `#161821` (Slate Steel) - For cards and containers.
- **Accents:** 
    - `#F43F5E` (Safety Red) - For alerts and stockouts.
    - `#38BDF8` (Blueprint Blue) - For dimensions and technical data.
    - `rgba(255, 255, 255, 0.05)` (Glass) - For subtle overlays.

## 2. Typography: "The Technical Standard"
- **Primary Font:** `Inter` (Sans-serif) - For readability and clean UI.
- **Monospace Font:** `JetBrains Mono` - For SKUs, Dimensions (2x4x8), and Prices. This reinforces the "Technical/Precision" feel.

## 3. The "WOW" Factors (Micro-Animations)
- **Fluid Transitions:** Page changes should use a subtle "Lens Blur" or "Fade-and-Scale" (0.98 -> 1.0) transition.
- **The "Pulse" of the Yard:** Real-time stock updates should have a subtle green glow pulse on the number when it changes.
- **Hover Depth:** Buttons and cards use `translateY(-2px)` with a soft glow shadow on hover to feel tactile.
- **Skeleton Loading:** Data loading should feel like a "Blueprint being drawn" rather than a spinner.

## 4. Layout: "Information Density without Clutter"
- **The "Command Bar" (Cmd+K):** A central search/action hub.
- **Contextual sidebars:** No deep nesting; the data comes to the user.
- **Dense Data-Grids:** High-performance grids with sticky headers and column pinning.

---

## Technical Tokens (CSS Variables preview)
```css
:root {
  --gable-green: #00FFA3;
  --gable-bg: #0A0B10;
  --gable-surface: #161821;
  --gable-text-main: #E2E8F0;
  --gable-font-ui: 'Inter', system-ui;
  --gable-font-data: 'JetBrains Mono', monospace;
  --gable-glass: rgba(255, 255, 255, 0.05);
  --gable-blur: blur(12px);
}
```
