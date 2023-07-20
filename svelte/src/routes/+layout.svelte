<script>
	import '../app.css';
	import { Header } from '$lib/components';
	import { onMount } from 'svelte';
	import { drawerToggled } from '$lib/store';

	onMount(async () => {
		checkWindowSize();
		window.addEventListener('resize', checkWindowSize);
	});

	function checkWindowSize() {
		if (window.matchMedia('(min-width: 768px)').matches) {
			if ($drawerToggled === false) {
				$drawerToggled = true;
			}
		} else {
			if ($drawerToggled === true) {
				$drawerToggled = false;
			}
		}
		// console.log('drawerToggled: ', $drawerToggled);
	}
</script>

<main>
	<Header />
	<div class="w-full h-16 min-h-[4rem]" />
	<slot />
</main>

<style lang="postcss">
	main {
		@apply flex flex-col w-full h-full;
	}
</style>
