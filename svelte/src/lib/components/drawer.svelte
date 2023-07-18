<script>
	import { drawerToggled } from '$lib/store';
	import Sidebar from './sidebar.svelte';
</script>

<div class="drawer-wrapper">
	<button
		class="drawer-mask {$drawerToggled ? '' : 'drawer-closed'}"
		on:click={() => ($drawerToggled = false)}
	/>
	<div class="drawer-content {$drawerToggled ? '' : 'drawer-closed'}">
		<Sidebar />
	</div>
	<div class="wrapper {$drawerToggled ? '' : 'drawer-closed'}">
		<slot />
	</div>
</div>

<style lang="postcss">
	.drawer-wrapper {
		@apply w-full h-full flex-grow relative overflow-hidden;
	}

	.drawer-mask {
		@apply w-full h-full bg-black/60 left-0 top-0 z-10 absolute;
		@apply transition-all duration-300 ease-in-out pointer-events-auto;
		@apply sm:hidden;
	}

	.drawer-closed.drawer-mask {
		@apply bg-transparent pointer-events-none;
	}

	.drawer-content {
		@apply w-48 h-full bg-white flex flex-col top-0 absolute !pointer-events-auto;
		@apply transition-[left] duration-300 ease-in-out left-0;
	}

	.drawer-closed.drawer-content {
		@apply -left-56;
	}

	.wrapper {
		@apply w-full h-full left-0 top-0 absolute z-0;
		@apply transition-[padding] duration-300 ease-in-out pl-0;
		@apply sm:pl-48;
	}

	.wrapper.drawer-closed {
		@apply pl-0;
	}
</style>
