<template>
  <div class="circular-progress">
    <svg viewBox="0 0 36 36" class="circular-chart">
      <path
        class="circle-bg"
        d="M18 2.0845
             a 15.9155 15.9155 0 0 1 0 31.831
             a 15.9155 15.9155 0 0 1 0 -31.831"
      />
      <path
        class="circle"
        :stroke-dasharray="`${percentage}, 100`"
        d="M18 2.0845
             a 15.9155 15.9155 0 0 1 0 31.831
             a 15.9155 15.9155 0 0 1 0 -31.831"
      />
      <text x="18" y="18" class="percentage">{{ percentage }}%</text>
      <text x="18" y="25" class="label">{{ label }}</text>
    </svg>
  </div>
</template>

<script setup>
defineProps({
  percentage: Number,
  label: String,
});
</script>

<style scoped>
.circular-progress {
  width: 80px;
  height: 80px;
}

.circular-chart {
  display: block;
  /* margin: 10px auto; */
  max-width: 100%;
  max-height: 100%;
}

.circle-bg {
  fill: none;
  stroke: #eee;
  stroke-width: 3.8;
}

.circle {
  fill: none;
  stroke-width: 2.8;
  stroke-linecap: round;
  stroke: #007bff;
  transition: stroke-dasharray 0.3s ease;
}

.percentage {
  fill: #007bff;
  font-family: Arial, sans-serif;
  font-size: 0.5em;
  text-anchor: middle;
  /* font-variant-caps: middle; */
}

.label {
  fill: #727272;
  font-family: Arial, sans-serif;
  font-size: 0.3em;
  text-anchor: middle;
  /* font-variant-caps: middle; */
}
</style>
