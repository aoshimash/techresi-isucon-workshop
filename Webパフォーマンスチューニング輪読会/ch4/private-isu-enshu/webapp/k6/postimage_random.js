import http from 'k6/http';

import { parseHTML } from "k6/html";

import { url } from "./config.js";

import { getAccount } from "./accounts.js";

const testImage = open("testimage.jpg", "b");

export default function() {

    const account = getAccount();
    
    const res = http.post(url("/login"), {
	account_name: account.account_name,
	password: account.password,
    });

    const doc = parseHTML(res.body);
    const token = doc.find('input[name="csrf_token"]').first().attr("value");

    http.post(url("/"), {
	file: http.file(testImage, "testImage.jpg", "image/jpeg"),
	body: "Posted by k6",
	csrf_token: token
    });
}
