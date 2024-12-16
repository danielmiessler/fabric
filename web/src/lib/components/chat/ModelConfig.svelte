<script lang="ts">
  import { Label } from "$lib/components/ui/label";
  import { Slider } from "$lib/components/ui/slider";
  import { modelConfig } from "$lib/store/model-store";
  import Transcripts from "./Transcripts.svelte";
</script>

    <div class="space-y-1">
        <Label>Maximum Length ({$modelConfig.maxLength})</Label>
        <Slider
            bind:value={$modelConfig.maxLength}
            min={1}
            max={4000}
            step={1}
        />
    </div>

    <div class="space-y-1">
        <Label>Temperature ({$modelConfig.temperature.toFixed(1)})</Label>
        <Slider
            bind:value={$modelConfig.temperature}
            min={0}
            max={2}
            step={0.1}
        />
    </div>

    <div class="space-y-1">
        <Label>Top P ({$modelConfig.top_p.toFixed(2)})</Label>
        <Slider
            bind:value={$modelConfig.top_p}
            min={0}
            max={1}
            step={0.01}
        />
    </div>

    <div class="space-y-1">
        <Label>Frequency Penalty ({$modelConfig.frequency.toFixed(2)})</Label>
        <Slider
            bind:value={$modelConfig.frequency}
            min={0}
            max={1}
            step={0.01}
        />
    </div>

    <div class="space-y-1">
        <Label>Presence Penalty ({$modelConfig.presence.toFixed(2)})</Label>
        <Slider
            bind:value={$modelConfig.presence}
            min={0}
            max={1}
            step={0.01}
        />
    </div>

    <div class="space-y-1">
        <Transcripts />
    </div>
