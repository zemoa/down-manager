import { writable, type Writable } from "svelte/store";
import type { Link } from "../model/link-model";

class LinkStore {
    constructor(
        public links: Writable<Link[]> = writable([]),
        public fetchingLinks: Writable<boolean> = writable(false)
    ) {}
}

export const linkStore = new LinkStore()