<script>
	import { cn } from "$lib/types/utils.ts";
	
	let className = undefined;
	export let value = 0;
	export let min = 0;
	export let max = 100;
	export let step = 1;
	export { className as class };

	let sliderEl;
	let isDragging = false;
	
	$: percentage = ((value - min) / (max - min)) * 100;

	function handleMouseDown(e) {
		isDragging = true;
		updateValue(e);
		window.addEventListener('mousemove', handleMouseMove);
		window.addEventListener('mouseup', handleMouseUp);
	}

	function handleMouseMove(e) {
		if (!isDragging) return;
		updateValue(e);
	}

	function handleMouseUp() {
		isDragging = false;
		window.removeEventListener('mousemove', handleMouseMove);
		window.removeEventListener('mouseup', handleMouseUp);
	}

	function updateValue(e) {
		if (!sliderEl) return;
		const rect = sliderEl.getBoundingClientRect();
		const pos = (e.clientX - rect.left) / rect.width;
		const rawValue = min + (max - min) * pos;
		const steppedValue = Math.round(rawValue / step) * step;
		value = Math.max(min, Math.min(max, steppedValue));
	}

	function handleKeyDown(e) {
		const step = e.shiftKey ? 10 : 1;
		switch (e.key) {
			case 'ArrowLeft':
			case 'ArrowDown':
				e.preventDefault();
				value = Math.max(min, value - step);
				break;
			case 'ArrowRight':
			case 'ArrowUp':
				e.preventDefault();
				value = Math.min(max, value + step);
				break;
			case 'Home':
				e.preventDefault();
				value = min;
				break;
			case 'End':
				e.preventDefault();
				value = max;
				break;
		}
	}
</script>

<div
	bind:this={sliderEl}
	class={cn("relative flex w-full touch-none select-none items-center", className)}
	on:mousedown={handleMouseDown}
	on:keydown={handleKeyDown}
	role="slider"
	tabindex="0"
	aria-valuemin={min}
	aria-valuemax={max}
	aria-valuenow={value}
>
	<div class="relative h-1.5 w-full grow overflow-hidden rounded-full bg-primary/20">
		<div 
			class="absolute h-full bg-primary transition-all" 
			style="width: {percentage}%" 
		/>
	</div>
	<div
		class="absolute h-4 w-4 rounded-full border border-secondary bg-primary-500 shadow transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50"
		style="left: calc({percentage}% - 0.5rem)"
	/>
</div>