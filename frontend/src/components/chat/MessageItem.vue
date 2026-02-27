<script setup lang="ts">
import { computed } from 'vue'
import ThinkingPanel from './ThinkingPanel.vue'
import { TerminalRectangleIcon, SecuredIcon, ThunderIcon, RobotIcon, UserIcon } from 'tdesign-icons-vue-next'

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
      let remaining = part
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
  <div class="flex gap-4 p-4 transition-colors rounded-lg hover:bg-bg-hover/50" :class="{'bg-bg-secondary': message.role === 'assistant'}">
    <div class="w-8 h-8 rounded-full flex items-center justify-center shrink-0 shadow-sm" 
      :class="message.role === 'user' ? 'bg-primary' : 'bg-success'">
      <UserIcon v-if="message.role === 'user'" class="text-white w-5 h-5" />
      <RobotIcon v-else class="text-white w-5 h-5" />
    </div>
    
    <div class="flex-1 min-w-0 space-y-2">
      <div class="font-semibold text-sm flex items-center gap-2 text-text-primary">
        {{ message.role === 'user' ? 'You' : 'FnChatBot' }}
        <span v-if="message.role === 'assistant'" class="text-xs font-normal text-text-muted bg-bg-primary px-1.5 py-0.5 rounded border border-border">
           {{ message.tasks?.length ? 'Planning' : 'Responding' }}
        </span>
      </div>
      
      <!-- Thinking Panel -->
      <ThinkingPanel 
        v-if="message.role === 'assistant' && (message.thinking || (message.tasks && message.tasks.length))"
        :thinking="message.thinking" 
        :tasks="message.tasks || []" 
      />

      <div class="space-y-2 text-sm text-text-primary">
        <template v-for="(seg, idx) in segments" :key="idx">
          
          <!-- Text -->
          <div v-if="seg.type === 'text'" class="whitespace-pre-wrap prose dark:prose-invert max-w-none">
            {{ seg.content }}
          </div>
          
          <!-- Bash (Sandbox) -->
          <div v-else-if="seg.type === 'bash'" class="relative group">
            <div class="absolute -top-3 left-2 bg-warning/20 text-warning text-[10px] font-bold px-2 py-0.5 rounded-full flex items-center gap-1 border border-warning/30 shadow-sm z-10">
              <SecuredIcon class="w-3 h-3" />
              Sandbox Mode
            </div>
            <div class="bg-black text-success p-4 rounded-lg font-mono text-xs overflow-x-auto border border-border mt-2">
              <div class="flex justify-between items-center mb-2 border-b border-gray-800 pb-2 opacity-50">
                <span class="flex items-center gap-1"><TerminalRectangleIcon class="w-3 h-3"/> bash</span>
              </div>
              <pre>{{ seg.content }}</pre>
            </div>
          </div>
          
          <!-- Other Code -->
          <div v-else-if="seg.type === 'code'" class="bg-bg-secondary p-4 rounded-lg font-mono text-xs overflow-x-auto border border-border">
            <div class="text-xs text-text-muted mb-1 opacity-70">{{ seg.meta }}</div>
            <pre>{{ seg.content }}</pre>
          </div>
          
          <!-- Subagent -->
          <div v-else-if="seg.type === 'subagent'" class="bg-info/10 border border-info/30 rounded-lg p-3 flex items-center gap-3">
             <div class="bg-info/20 p-1.5 rounded-md text-info">
               <RobotIcon class="w-4 h-4" />
             </div>
             <div>
               <div class="text-xs font-bold text-info uppercase tracking-wider">Subagent: {{ seg.meta }}</div>
               <div class="text-text-primary">{{ seg.content }}</div>
             </div>
          </div>
          
          <!-- Skill -->
          <div v-else-if="seg.type === 'skill'" class="bg-accent/10 border border-accent/30 rounded-lg p-3 flex items-center gap-3">
             <div class="bg-accent/20 p-1.5 rounded-md text-accent">
               <ThunderIcon class="w-4 h-4" />
             </div>
             <div>
               <div class="text-xs font-bold text-accent uppercase tracking-wider">Skill Loaded</div>
               <div class="text-text-primary">{{ seg.content }}</div>
             </div>
          </div>

        </template>
      </div>
    </div>
  </div>
</template>

