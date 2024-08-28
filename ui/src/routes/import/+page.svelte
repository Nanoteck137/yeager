<script lang="ts">
  import { enhance } from "$app/forms";

  let files = $state<File[]>([]);

  let fileSelector: HTMLInputElement;
</script>

<p>Import an Album</p>

<form
  class="flex flex-col gap-2"
  method="POST"
  enctype="multipart/form-data"
  use:enhance={({ formData }) => {
    files.forEach((f) => {
      formData.append("files", f);
    });
  }}
>
  <input
    class="border-gray-600 bg-gray-700 text-white placeholder:text-white"
    name="albumName"
    placeholder="Album Name"
    type="text"
    minlength="1"
  />
  <input
    class="border-gray-600 bg-gray-700 text-white placeholder:text-white"
    name="artistName"
    placeholder="Artist Name"
    type="text"
  />

  {#each files as file}
    <p>{file.name}</p>
  {/each}

  <button
    type="button"
    class="w-fit rounded bg-purple-300 px-6 py-2 hover:bg-purple-100 hover:text-black active:scale-95"
    onclick={() => {
      fileSelector.click();
    }}>Add Files</button
  >

  <button
    class="w-fit rounded bg-blue-500 px-6 py-2 hover:bg-blue-400 active:scale-95"
    >Import</button
  >
</form>

<input
  class="hidden"
  bind:this={fileSelector}
  multiple
  type="file"
  accept="audio/*"
  onchange={(e) => {
    e.preventDefault();

    const target = e.target as HTMLInputElement;
    if (!target.files) {
      return;
    }

    for (let i = 0; i < target.files.length; i++) {
      const item = target.files[i];
      if (!files.find((f) => f.name === item.name)) {
        files.push(item);
      }
    }
  }}
/>
