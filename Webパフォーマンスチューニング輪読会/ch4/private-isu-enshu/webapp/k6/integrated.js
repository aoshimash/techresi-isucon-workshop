import initialize from "./initialize.js";
import comment_random from "./comment_random.js";
import postimage_random from "./postimage_random.js";

export { initialize, comment_random, postimage_random };

export const options = {
    scenarios: {
	initialize: {
	    executor: "shared-iterations",
	    vus: 1,
	    iterations: 1,
	    exec: "initialize",
	    maxDuration: "10s"
	},
	comment: {
	    executor: "constant-vus",
	    vus: 4,
	    duration: "30s",
	    exec: "comment_random",
	    startTime: "12s",
	},
	postImage: {
	    executor: "constant-vus",
	    vus: 2,
	    duration: "30s",
	    exec: "postimage_random",
	    startTime: "12s",
	},
    },
};

export default function () { }
