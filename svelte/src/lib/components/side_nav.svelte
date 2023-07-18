<script>
	import { drawerToggled } from '$lib/store';
</script>

<div class="sidenav-container">
	<div class="sidenav-content {$drawerToggled === false ? 'sidenav-hidden' : ''}">
		<slot />
	</div>
</div>

<style lang="postcss">
	.sidenav-container {
		@apply h-full relative;
	}
	.sidenav-content {
		@apply block h-full transition-[margin] duration-300;
		@apply md:ml-64 lg:ml-80;
	}
	.sidenav-content.sidenav-hidden {
		@apply ml-0;
	}

	.sidenav-content {
		@apply block relative h-full;
	}
</style>
