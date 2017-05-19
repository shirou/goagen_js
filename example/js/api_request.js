// This module exports functions that give access to the goagen_js API hosted at localhost:8080.

import 'whatwg-fetch';

import * as v from "./api_validator.js";

const scheme = 'http';
const host = 'localhost:8080';
const urlPrefix = scheme + '://' + host;

// GetInt
// Get Method
// 
export function GetGetInt(payload) {
  const url = urlPrefix + "/get/int";
  let errors = {};
  let ret;
  ret = v.validate(v.GetGetInt.payload, payload);
  if (ret !== null) {
     errors.payload = ret;
  }
  if (errors.length > 0) {
    return Promise.reject({status: 400, detail: "validatoin error", meta: errors});
  }
  return get(url, payload);
}

// PathParams
// Get Method with params in path
// 
// paramInt(number): path_params param int
// paramStr(string): 
export function GetPathParams(paramInt, paramStr) {
  const url = urlPrefix + "/get/int/${paramInt}/${paramStr}";
  let errors = {};
  let ret;
  ret = v.validate(v.GetPathParams.paramInt, paramInt);
  if (ret !== null) {
     errors.paramInt = ret;
  }
  ret = v.validate(v.GetPathParams.paramStr, paramStr);
  if (ret !== null) {
     errors.paramStr = ret;
  }
  if (errors.length > 0) {
    return Promise.reject({status: 400, detail: "validatoin error", meta: errors});
  }
  return get(url);
}

// Without
// Get Method without params
// 
export function GetWithout() {
  const url = urlPrefix + "/get";
  return get(url);
}

// helper function for GET method.
function get(url, payload) {
  const query = queryBuilder(payload);
  return fetch(url + query, {
    method: 'GET',
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
    },
  });
}

// helper function for POST method.
function post(url, payload) {
  return fetch(url, {
    method: 'POST',
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  });
}

// helper functon which return QueryParameter from Object.
function queryBuilder(obj) {
  if (!obj) {
    return '';
  }
  const r = Object.keys(obj).sort().map((key) => {
    const val = obj[key];
    if (val === undefined){
      return '';
    }
    if (val === null){
      return '';
    }
    return encodeURIComponent(key) + "=" + encodeURIComponent(val);
  }).filter((x) => {
    return x.length > 0;
  }).join('&');
  if (r.length > 0){
    return '?' + r;
  }
  return '';
}
