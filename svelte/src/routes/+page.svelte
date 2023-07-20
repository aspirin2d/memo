<script>
	import SideNav from '../lib/components/side_nav.svelte';
	import ChatWindow from '../lib/components/chat-window.svelte';
</script>

<SideNav>
	<div slot="drawer">Thiis is DRAWER</div>
	<div slot="content">
		<ChatWindow />
	</div>
</SideNav>
