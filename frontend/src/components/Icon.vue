<template>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    :width="size"
    :height="size"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="1.8"
    stroke-linecap="round"
    stroke-linejoin="round"
  >
    <template v-for="(d, i) in paths" :key="i">
      <path v-if="d.type === 'path'" :d="d.d" />
      <circle v-else-if="d.type === 'circle'" :cx="d.cx" :cy="d.cy" :r="d.r" />
      <rect v-else-if="d.type === 'rect'" :x="d.x" :y="d.y" :width="d.w" :height="d.h" :rx="d.rx" />
      <line v-else-if="d.type === 'line'" :x1="d.x1" :y1="d.y1" :x2="d.x2" :y2="d.y2" />
      <polyline v-else-if="d.type === 'polyline'" :points="d.points" />
    </template>
  </svg>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  name: { type: String, required: true },
  size: { type: [Number, String], default: 20 },
})

// Lightweight inline icon set (Lucide-style geometry).
const ICONS = {
  dashboard: [
    { type: 'rect', x: 3, y: 3, w: 7, h: 9, rx: 1 },
    { type: 'rect', x: 14, y: 3, w: 7, h: 5, rx: 1 },
    { type: 'rect', x: 14, y: 12, w: 7, h: 9, rx: 1 },
    { type: 'rect', x: 3, y: 16, w: 7, h: 5, rx: 1 },
  ],
  users: [
    { type: 'path', d: 'M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2' },
    { type: 'circle', cx: 9, cy: 7, r: 4 },
    { type: 'path', d: 'M22 21v-2a4 4 0 0 0-3-3.87' },
    { type: 'path', d: 'M16 3.13a4 4 0 0 1 0 7.75' },
  ],
  list: [
    { type: 'line', x1: 8, y1: 6, x2: 21, y2: 6 },
    { type: 'line', x1: 8, y1: 12, x2: 21, y2: 12 },
    { type: 'line', x1: 8, y1: 18, x2: 21, y2: 18 },
    { type: 'line', x1: 3, y1: 6, x2: 3.01, y2: 6 },
    { type: 'line', x1: 3, y1: 12, x2: 3.01, y2: 12 },
    { type: 'line', x1: 3, y1: 18, x2: 3.01, y2: 18 },
  ],
  plus: [
    { type: 'line', x1: 12, y1: 5, x2: 12, y2: 19 },
    { type: 'line', x1: 5, y1: 12, x2: 19, y2: 12 },
  ],
  settings: [
    { type: 'circle', cx: 12, cy: 12, r: 3 },
    { type: 'path', d: 'M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z' },
  ],
  chevron: [{ type: 'polyline', points: '9 18 15 12 9 6' }],
  menu: [
    { type: 'line', x1: 3, y1: 6, x2: 21, y2: 6 },
    { type: 'line', x1: 3, y1: 12, x2: 21, y2: 12 },
    { type: 'line', x1: 3, y1: 18, x2: 21, y2: 18 },
  ],
  bell: [
    { type: 'path', d: 'M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9' },
    { type: 'path', d: 'M13.73 21a2 2 0 0 1-3.46 0' },
  ],
  search: [
    { type: 'circle', cx: 11, cy: 11, r: 8 },
    { type: 'line', x1: 21, y1: 21, x2: 16.65, y2: 16.65 },
  ],
  logout: [
    { type: 'path', d: 'M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4' },
    { type: 'polyline', points: '16 17 21 12 16 7' },
    { type: 'line', x1: 21, y1: 12, x2: 9, y2: 12 },
  ],
  edit: [
    { type: 'path', d: 'M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7' },
    { type: 'path', d: 'M18.5 2.5a2.12 2.12 0 0 1 3 3L12 15l-4 1 1-4z' },
  ],
  trash: [
    { type: 'polyline', points: '3 6 5 6 21 6' },
    { type: 'path', d: 'M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2' },
  ],
  shield: [{ type: 'path', d: 'M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z' }],
  check: [{ type: 'polyline', points: '20 6 9 17 4 12' }],
  x: [
    { type: 'line', x1: 18, y1: 6, x2: 6, y2: 18 },
    { type: 'line', x1: 6, y1: 6, x2: 18, y2: 18 },
  ],
  user: [
    { type: 'path', d: 'M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2' },
    { type: 'circle', cx: 12, cy: 7, r: 4 },
  ],
  activity: [{ type: 'polyline', points: '22 12 18 12 15 21 9 3 6 12 2 12' }],
  arrowLeft: [
    { type: 'line', x1: 19, y1: 12, x2: 5, y2: 12 },
    { type: 'polyline', points: '12 19 5 12 12 5' },
  ],
  building: [
    { type: 'rect', x: 4, y: 2, w: 16, h: 20, rx: 2 },
    { type: 'line', x1: 9, y1: 22, x2: 9, y2: 18 },
    { type: 'line', x1: 15, y1: 22, x2: 15, y2: 18 },
    { type: 'line', x1: 8, y1: 7, x2: 8.01, y2: 7 },
    { type: 'line', x1: 12, y1: 7, x2: 12.01, y2: 7 },
    { type: 'line', x1: 16, y1: 7, x2: 16.01, y2: 7 },
    { type: 'line', x1: 8, y1: 12, x2: 8.01, y2: 12 },
    { type: 'line', x1: 12, y1: 12, x2: 12.01, y2: 12 },
    { type: 'line', x1: 16, y1: 12, x2: 16.01, y2: 12 },
  ],
  briefcase: [
    { type: 'rect', x: 2, y: 7, w: 20, h: 14, rx: 2 },
    { type: 'path', d: 'M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16' },
  ],
  calendar: [
    { type: 'rect', x: 3, y: 4, w: 18, h: 18, rx: 2 },
    { type: 'line', x1: 16, y1: 2, x2: 16, y2: 6 },
    { type: 'line', x1: 8, y1: 2, x2: 8, y2: 6 },
    { type: 'line', x1: 3, y1: 10, x2: 21, y2: 10 },
  ],
}

const paths = computed(() => ICONS[props.name] || [])
</script>
