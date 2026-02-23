<script setup lang="ts">
import { computed } from 'vue'
import ThinkingPanel from './ThinkingPanel.vue'
import { Terminal, ShieldCheck, Zap, Bot } from 'lucide-vue-next'

const props = defineProps<{
  message: any
}>()

interface Segment {
  type: 'text' | 'bash' | 'code' | 'subagent' | 'skill'
  content: string
  meta?: string
}

const segments = computed(() => {
  const content = props.message.content || ''
  const result: Segment[] = []
  
  // Split by code blocks first
  const parts = content.split(/(```[\s\S]*?```)/g)
  
  parts.forEach(part => {
    if (part.startsWith('```')) {
      const match = part.match(/```(\w+)?\n([\s\S]*?)```/)
      if (match) {
        const lang = match[1] || ''
        const code = match[2]
        if (lang === 'bash' || lang === 'sh') {
          result.push({ type: 'bash', content: code, meta: lang })
        } else {
          result.push({ type: 'code', content: code, meta: lang })
        }
      } else {
        result.push({ type: 'text', content: part })
      }
    } else {
      // Process text for Subagent/Skill patterns
      // Patterns from backend: 
      // *Subagent [Type] started: Desc*
      // *Loaded Skill: Name*
      
      let remaining = part
      
      // We need a loop to find all occurrences
      // Regex with capture groups
      const subagentRegex = /\*Subagent \[(.*?)\] started: (.*?)\*/g
      const skillRegex = /\*Loaded Skill: (.*?)\*/g
      
      // Simplification: Just split by lines and check patterns if regex split is complex
      // Or use a tokenization approach.
      // Let's just process lines for simplicity in this demo scope
      
      const lines = remaining.split('\n')
      let currentText = ''
      
      lines.forEach(line => {
        const subagentMatch = line.match(/^\*Subagent \[(.*?)\] started: (.*?)\*$/)
        const skillMatch = line.match(/^\*Loaded Skill: (.*?)\*$/)
        
        if (subagentMatch) {
          if (currentText) {
             result.push({ type: 'text', content: currentText })
             currentText = ''
          }
          result.push({ type: 'subagent', content: subagentMatch[2], meta: subagentMatch[1] })
        } else if (skillMatch) {
          if (currentText) {
             result.push({ type: 'text', content: currentText })
             currentText = ''
          }
          result.push({ type: 'skill', content: skillMatch[1] })
        } else {
          currentText += line + '\n'
        }
      })
      
      if (currentText) {
        result.push({ type: 'text', content: currentText })
      }
    }
  })
  
  return result
})
</script>

<template>
  <div class="flex gap-4 p-4 transition-colors" :class="{'bg-zinc-50 dark:bg-zinc-900/50': message.role === 'assistant'}">
    <div class="w-8 h-8 rounded-full flex items-center justify-center shrink-0 shadow-sm" 
      :class="message.role === 'user' ? 'bg-blue-500' : 'bg-green-600'">
      <span class="text-white text-xs font-bold">{{ message.role === 'user' ? 'U' : 'AI' }}</span>
    </div>
    
    <div class="flex-1 min-w-0 space-y-2">
      <div class="font-semibold text-sm flex items-center gap-2">
        {{ message.role === 'user' ? 'You' : 'FnChatBot' }}
        <span v-if="message.role === 'assistant'" class="text-xs font-normal text-muted-foreground bg-secondary px-1.5 py-0.5 rounded">
           {{ message.tasks?.length ? 'Planning' : 'Responding' }}
        </span>
      </div>
      
      <!-- Thinking Panel -->
      <ThinkingPanel 
        v-if="message.role === 'assistant' && (message.thinking || (message.tasks && message.tasks.length))"
        :thinking="message.thinking" 
        :tasks="message.tasks || []" 
      />

      <div class="space-y-2 text-sm">
        <template v-for="(seg, idx) in segments" :key="idx">
          
          <!-- Text -->
          <div v-if="seg.type === 'text'" class="whitespace-pre-wrap prose dark:prose-invert max-w-none">
            {{ seg.content }}
          </div>
          
          <!-- Bash (Sandbox) -->
          <div v-else-if="seg.type === 'bash'" class="relative group">
            <div class="absolute -top-3 left-2 bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-400 text-[10px] font-bold px-2 py-0.5 rounded-full flex items-center gap-1 border border-amber-200 dark:border-amber-800 shadow-sm z-10">
              <ShieldCheck class="w-3 h-3" />
              Sandbox Mode
            </div>
            <div class="bg-black text-green-400 p-4 rounded-lg font-mono text-xs overflow-x-auto border border-zinc-800 mt-2">
              <div class="flex justify-between items-center mb-2 border-b border-zinc-800 pb-2 opacity-50">
                <span class="flex items-center gap-1"><Terminal class="w-3 h-3"/> bash</span>
              </div>
              <pre>{{ seg.content }}</pre>
            </div>
          </div>
          
          <!-- Other Code -->
          <div v-else-if="seg.type === 'code'" class="bg-muted p-4 rounded-lg font-mono text-xs overflow-x-auto border border-border">
            <div class="text-xs text-muted-foreground mb-1 opacity-70">{{ seg.meta }}</div>
            <pre>{{ seg.content }}</pre>
          </div>
          
          <!-- Subagent -->
          <div v-else-if="seg.type === 'subagent'" class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-3 flex items-center gap-3">
             <div class="bg-blue-100 dark:bg-blue-800 p-1.5 rounded-md text-blue-600 dark:text-blue-300">
               <Bot class="w-4 h-4" />
             </div>
             <div>
               <div class="text-xs font-bold text-blue-700 dark:text-blue-300 uppercase tracking-wider">Subagent: {{ seg.meta }}</div>
               <div class="text-blue-600 dark:text-blue-200">{{ seg.content }}</div>
             </div>
          </div>
          
          <!-- Skill -->
          <div v-else-if="seg.type === 'skill'" class="bg-purple-50 dark:bg-purple-900/20 border border-purple-200 dark:border-purple-800 rounded-lg p-3 flex items-center gap-3">
             <div class="bg-purple-100 dark:bg-purple-800 p-1.5 rounded-md text-purple-600 dark:text-purple-300">
               <Zap class="w-4 h-4" />
             </div>
             <div>
               <div class="text-xs font-bold text-purple-700 dark:text-purple-300 uppercase tracking-wider">Skill Loaded</div>
               <div class="text-purple-600 dark:text-purple-200">{{ seg.content }}</div>
             </div>
          </div>

        </template>
      </div>
    </div>
  </div>
</template>
