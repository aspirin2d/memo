<script>
	import { fade } from 'svelte/transition';
	import { onMount } from 'svelte';

	var commands = [
		{ label: '设置启动模版', icon: 'T' },
		{ label: '设置变量', icon: 'V' }
	];

	onMount(() => {
		console.log('mounted');
	});
</script>

<div class="command-list" transition:fade={{ duration: 250 }}>
	{#each commands as command}
		<button class="command">
			<span class="kbd">{command.icon}</span>
			<span>{command.label}</span>
		</button>
	{/each}
</div>

<style lang="postcss">
	.command-list {
		@apply absolute bottom-[56px] bg-white py-2 px-4 min-w-[256px] rounded-3xl shadow
  mb-4 space-y-2;
	}

	.command {
		@apply flex flex-row items-center h-8 w-full;
	}

	.kbd {
		@apply border border-gray-400 border-b-2 min-w-[2.2em] min-h-[2.2em] flex items-center
  justify-center mr-2 p-0.5 rounded-md text-sm;
	}
</style>
