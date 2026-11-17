export const googleError: Record<string, string> = {
  google: 'Google sign-in did not finish. Try again.',
  google_email: 'Google must use the same email as this account.',
  google_taken: 'That Google account is already used on another login.',
  google_unverified: 'Google has not confirmed that email.',
}

export function googleErrorMessage(code: unknown) {
  if (typeof code !== 'string' || code === '') return ''
  return googleError[code] ?? googleError.google
}
