<script>
	import { drawerToggled } from '$lib/store';
	import { fly, fade } from 'svelte/transition';
</script>

<div class="sidenav-container {$drawerToggled === false ? 'sidenav-hidden' : ''}">
	<div class="sidenav-drawer">
		{#if $drawerToggled}
			<button
				class="sidenav-drawer-backdrop"
				transition:fade
				on:click={() => {
					$drawerToggled = false;
				}}
			/>
			<div class="sidenav-drawer-content" transition:fly={{ x: '-100%' }}>
				<slot name="drawer" />
			</div>
		{/if}
	</div>
	<div class="sidenav-content">
		<slot name="content" />
	</div>
</div>

<style lang="postcss">
	.sidenav-drawer {
		@apply absolute top-0 bottom-0 overflow-hidden overflow-y-auto
  z-10 w-full sm:w-64 lg:w-80;
	}

	.sidenav-drawer-backdrop {
		@apply bg-black/60 absolute top-0 bottom-0 left-0 right-0 z-20 sm:hidden;
	}

	.sidenav-drawer-content {
		@apply absolute top-0 bottom-0 left-0 right-0 z-30 mr-16 sm:mr-0;
	}

	.sidenav-container {
		@apply h-full relative;
	}
	.sidenav-content {
		@apply block h-full transition-[margin] duration-[400ms];
		@apply sm:ml-64 lg:ml-80;
	}
	.sidenav-hidden .sidenav-content {
		@apply ml-0;
	}

	.sidenav-content {
		@apply block relative h-full;
	}

	:global([slot='content']) {
		@apply h-full;
	}
</style>
