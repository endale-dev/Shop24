<template>
  <div style="display:flex;gap:16px">
    <aside style="width:280px;border-right:1px solid #eee;padding:12px">
      <h3>Blocks</h3>
      <ul>
        <li v-for="b in blocks" :key="b.id">{{ b.componentKey }}</li>
      </ul>
    </aside>
    <main style="flex:1;padding:12px">
      <h2>Page Editor (WYSIWYG + Tree)</h2>
      <div>
        <label>Title</label>
        <input v-model="title" />
      </div>
      <div style="margin-top:12px">
        <editor-content :editor="editor" />
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { EditorContent, useEditor } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'

const title = ref('New Page')
const blocks = ref([{ id: 'root', componentKey: 'Container', props: {} }])

const editor = useEditor({ extensions: [StarterKit], content: '<p>Edit page content here</p>' })
</script>
