# Design System Documentation

## Overview
The sellersDashboard uses a comprehensive design system to ensure consistency across all UI components and pages.

## Colors

### Primary Palette
- **Blue (Primary)**: Used for primary actions, links, and active states
  - `500`: #3b82f6 (main)
  - `600`: #2563eb (hover)
  - `700`: #1d4ed8 (active)

### Semantic Colors
- **Success (Green)**: #22c55e - Positive actions, completed states
- **Warning (Amber)**: #f59e0b - Warnings, pending states  
- **Danger (Red)**: #ef4444 - Errors, destructive actions
- **Neutral (Gray)**: #6b7280 - Text, borders, backgrounds

## Typography

### Font Family
- **Primary**: Inter (sans-serif)
- **Monospace**: Monaco, Consolas

### Font Sizes
- `xs`: 0.75rem (12px)
- `sm`: 0.875rem (14px)
- `base`: 1rem (16px)
- `lg`: 1.125rem (18px)
- `xl`: 1.25rem (20px)
- `2xl`: 1.5rem (24px)
- `3xl`: 1.875rem (30px)

### Font Weights
- Normal: 400
- Medium: 500
- Semibold: 600
- Bold: 700

## Spacing Scale
- 1: 0.25rem (4px)
- 2: 0.5rem (8px)
- 3: 0.75rem (12px)
- 4: 1rem (16px)
- 6: 1.5rem (24px)
- 8: 2rem (32px)

## Animations

### Available Animations
- `animate-fade-in`: Fade in animation
- `animate-slide-up`: Slide up with fade
- `animate-slide-down`: Slide down with fade
- `animate-scale-in`: Scale with fade
- `animate-spin`: Loading spinner
- `animate-pulse`: Pulsing effect

### Usage
```vue
<div class="animate-slide-up">Content</div>
```

## Components

### Base Components
1. **BaseButton** - All button variants
2. **BaseCard** - Content containers
3. **BaseBadge** - Status indicators
4. **BaseInput** - Form inputs
5. **BaseAlert** - Notifications
6. **DataTable** - Data tables

### Common Components
1. **StatCard** - Metric cards
2. **PageHeader** - Page headers
3. **SkeletonLoader** - Loading states

## Responsive Design

### Breakpoints
- `sm`: 640px
- `md`: 768px
- `lg`: 1024px
- `xl`: 1280px

### Mobile-First Approach
All designs start mobile and progressively enhance for larger screens.

## Utilities

### Formatters (`utils/formatters.js`)
- `formatCurrency(amount)` - Format money values
- `formatDate(date)` - Format dates
- `truncate(text, length)` - Truncate long text

### Validation (`utils/validation.js`)
- `validationRules.required`
- `validationRules.email`
- `validationRules.min / max`
- `validationRules.price`

### Animations (`utils/animations.js`)
- `staggerDelay(index)` - Stagger list animations
- `useScrollAnimation()` - Scroll-triggered animations

## Best Practices

1. **Use design tokens** instead of hard-coded values
2. **Prefer utility classes** for common styles
3. **Use semantic color names** (primary, success, etc.)
4. **Follow mobile-first** responsive design
5. **Add transitions** to interactive elements
6. **Use proper accessibility** attributes
