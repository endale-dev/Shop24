<template>
  <div class="grid min-h-screen lg:grid-cols-2">
    <section class="hidden bg-ink px-12 py-14 text-stone-100 lg:flex lg:flex-col lg:justify-between">
      <p class="font-serif text-3xl">Sable</p>
      <div>
        <h1 class="max-w-md font-serif text-5xl leading-tight">Describe a shop. Publish the storefront the same day.</h1>
        <p class="mt-4 max-w-sm text-sm text-stone-300">A small commerce studio for one seller, one catalog, and a design that stays tied to real products.</p>
      </div>
      <p class="text-xs uppercase tracking-[0.2em] text-stone-400">Cash on delivery · No theme code</p>
    </section>
    <section class="flex items-center justify-center bg-paper px-5 py-12">
      <div class="w-full max-w-md">
        <p class="mb-6 font-serif text-3xl lg:hidden">Sable</p>
        <RouterView />
      </div>
    </section>
  </div>
</template>
