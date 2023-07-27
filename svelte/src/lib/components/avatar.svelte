<script>
	import { createAvatar } from '@dicebear/core';
	import { micah } from '@dicebear/collection';

	export let seed;
	const avatar = createAvatar(micah, {
		seed: seed
	}).toDataUriSync();
</script>

<div class="avatar">
	<img src={avatar} alt="Avatar" />
</div>

<style lang="postcss">
	.avatar {
		@apply w-8 h-8 rounded-full;
	}
</style>
