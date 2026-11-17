const controlBase =
  'w-full rounded border border-border bg-surface px-2.5 text-sm text-text outline-none placeholder:text-muted focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent/30'

/** Single-line inputs and selects. */
export const controlClass = `h-8 ${controlBase}`

/**
 * Textareas use the same chrome as controlClass without the fixed height.
 * Callers that stay on controlClass should add py-2 and drop h-8.
 */
export const textareaClass = `py-2 ${controlBase}`
