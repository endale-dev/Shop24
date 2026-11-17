<script setup lang="ts">
import Button from '@/components/ui/Button.vue'
import Logo from '@/components/ui/Logo.vue'
</script>

<template>
  <div class="min-h-screen bg-bg text-text">
    <header class="sticky top-0 z-10 border-b border-border bg-surface">
      <div class="mx-auto flex h-12 max-w-5xl items-center justify-between gap-3 px-4">
        <RouterLink to="/" class="inline-flex items-center gap-2 text-text">
          <Logo />
          <span class="text-sm font-medium leading-none">Sable</span>
        </RouterLink>
        <nav class="flex items-center gap-3 text-sm">
          <a class="hidden text-muted hover:text-text sm:inline" href="#shop">The shop</a>
          <a class="hidden text-muted hover:text-text sm:inline" href="#desk">The desk</a>
          <RouterLink class="text-text hover:underline" to="/login">Sign in</RouterLink>
          <Button to="/register">Open a shop</Button>
        </nav>
      </div>
    </header>

    <main>
      <section class="mx-auto grid max-w-5xl items-center gap-8 px-4 py-12 md:grid-cols-2 md:py-16">
        <div>
          <p class="text-sm text-muted">One seller. Cash when the parcel arrives.</p>
          <h1 class="mt-2 text-3xl font-semibold leading-tight">The desk for a shop you pack yourself.</h1>
          <p class="mt-4 text-sm leading-6 text-muted">
            Sable is for a single owner. You list the goods, the versions you actually sell, and the groups you want on the page. Customers open the shop on its own address, choose a color or a size, and pay the rider. There is no card form and no second seller in the account.
          </p>
          <div class="mt-6 flex flex-wrap items-center gap-3">
            <Button to="/register">Create the account</Button>
            <RouterLink class="text-sm font-medium text-text underline" to="/login">I already sell</RouterLink>
          </div>
        </div>
        <img
          class="aspect-[4/3] w-full rounded border border-border object-cover"
          src="/landing/rail.jpg"
          alt="A charcoal wool coat and a sand wool coat hanging on a wooden rail"
        />
      </section>

      <section id="day" class="border-t border-border bg-surface">
        <div class="mx-auto grid max-w-5xl gap-8 px-4 py-12 md:grid-cols-[1.1fr_0.9fr] md:items-start">
          <div>
            <h2 class="text-lg font-semibold leading-6">A day with the goods</h2>
            <ol class="mt-4 space-y-4 text-sm leading-6">
              <li>
                <span class="font-medium">List what you can hand over.</span>
                <span class="text-muted"> A product has a name, a price kept in cents, stock, and a picture address. A coat at $180 is stored as 18000 cents so the total does not drift.</span>
              </li>
              <li>
                <span class="font-medium">Add the versions you sell.</span>
                <span class="text-muted"> Color, size, capacity: the option names are yours. Charcoal at $180 with three in stock, Sand at $190 with one, each with its own picture and SKU. Sable does not invent the other combinations.</span>
              </li>
              <li>
                <span class="font-medium">Group them by hand.</span>
                <span class="text-muted"> A collection is a list you choose, with its own picture. Winter coats contains the coat. It does not fill itself from a rule.</span>
              </li>
              <li>
                <span class="font-medium">Take the order.</span>
                <span class="text-muted"> The customer leaves a name, a phone, and an address. You collect the cash when the parcel changes hands. Stock comes down on that order.</span>
              </li>
            </ol>
          </div>
          <img
            class="aspect-[4/3] w-full rounded border border-border object-cover"
            src="/landing/parcel.jpg"
            alt="A kraft paper parcel tied with cotton tape on a wooden table"
          />
        </div>
      </section>

      <section id="shop" class="mx-auto grid max-w-5xl items-center gap-8 px-4 py-12 md:grid-cols-[0.8fr_1.2fr]">
        <img
          class="aspect-[3/4] w-full rounded border border-border object-cover md:max-h-[520px]"
          src="/landing/shop.jpg"
          alt="Folded linen on a wooden shelf and a coat on a hook in a small shop"
        />
        <div>
          <h2 class="text-lg font-semibold leading-6">Two addresses, on purpose</h2>
          <p class="mt-3 text-sm leading-6 text-muted">
            This desk stays on the bare domain. The shop the customer sees lives on a subdomain of the same parent. Locally that is <span class="text-text">your-shop.localhost</span> beside this page. The seller session is not sent to the shop host, so a customer browser never carries the login.
          </p>
          <p class="mt-3 text-sm leading-6 text-muted">
            The same-host preview, <span class="text-text">/s/your-shop</span>, is for checking the page while you work. The subdomain is the address you give out.
          </p>
          <dl class="mt-6 grid gap-3 text-sm sm:grid-cols-2">
            <div class="rounded border border-border bg-surface p-3">
              <dt class="font-medium">Desk</dt>
              <dd class="mt-1 text-muted">Products, collections, orders, the builder, and the store settings. Overview is here even before the shop exists.</dd>
            </div>
            <div class="rounded border border-border bg-surface p-3">
              <dt class="font-medium">Shop</dt>
              <dd class="mt-1 text-muted">Home, a collection, a product, the cart, and checkout. The cart stays in the browser. Checkout is cash on delivery.</dd>
            </div>
          </dl>
        </div>
      </section>

      <section class="border-t border-border bg-surface">
        <div class="mx-auto max-w-5xl px-4 py-12">
          <h2 class="text-lg font-semibold leading-6">What the page is allowed to change</h2>
          <p class="mt-3 max-w-3xl text-sm leading-6 text-muted">
            The builder writes a layout the shop already knows how to draw. You can ask for the colors, the type, how tight the page is, whether buttons are filled or just links, how the hero sits, and whether products are a grid, a list, or one feature. It does not paste a page of code into the shop.
          </p>
          <div class="mt-6 grid gap-3 text-sm md:grid-cols-3">
            <div class="rounded border border-border p-3">
              <p class="font-medium">A prompt, then the page</p>
              <p class="mt-1 text-muted">Send a change and the draft updates. If the shop is already published, that same layout is what customers see. The next prompt starts from the conversation, not from a blank page.</p>
            </div>
            <div class="rounded border border-border p-3">
              <p class="font-medium">Sections you can hide</p>
              <p class="mt-1 text-muted">Hero, featured products, collections, about, an announcement, and a newsletter note. The newsletter is a line of copy. It does not collect email addresses.</p>
            </div>
            <div class="rounded border border-border p-3">
              <p class="font-medium">An honest stand-in</p>
              <p class="mt-1 text-muted">With no model configured, the built-in designer still fills a layout and says so. It does not pretend a live model answered.</p>
            </div>
          </div>
        </div>
      </section>

      <section id="desk" class="mx-auto grid max-w-5xl items-center gap-8 px-4 py-12 md:grid-cols-2">
        <div>
          <h2 class="text-lg font-semibold leading-6">The desk stays plain</h2>
          <p class="mt-3 text-sm leading-6 text-muted">
            After you sign in, the sidebar is the shop: overview, products, collections, orders, builder, settings. Collapse it when you want the builder to take the width. The account sits on the initials at the bottom of that sidebar, and on the small bar on a phone.
          </p>
          <ul class="mt-4 space-y-2 text-sm leading-6 text-muted">
            <li>Sign in with email and a password. Reset it from the sign-in page. Change it, or set one for the first time, from the account page.</li>
            <li>Google is there when this server has been given a client id and secret. Connecting it from the account page keeps the same email. It does not open a second shop.</li>
            <li>One account holds one store. Prices, stock, and orders stay on that store.</li>
          </ul>
          <div class="mt-6">
            <Button to="/register">Open a shop</Button>
          </div>
        </div>
        <img
          class="aspect-video w-full rounded border border-border object-cover"
          src="/landing/desk.jpg"
          alt="A wooden desk with a notebook, a ceramic cup, a tape measure, and a wool sleeve"
        />
      </section>
    </main>

    <footer class="border-t border-border">
      <div class="mx-auto flex max-w-5xl flex-wrap items-center justify-between gap-3 px-4 py-6 text-sm text-muted">
        <p>Sable. Cash on delivery, one seller.</p>
        <div class="flex gap-4">
          <RouterLink class="hover:text-text" to="/login">Sign in</RouterLink>
          <RouterLink class="hover:text-text" to="/register">Create account</RouterLink>
        </div>
      </div>
    </footer>
  </div>
</template>
