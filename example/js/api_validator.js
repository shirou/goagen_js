// This module exports functions that validate goagen_js API params hosted at localhost:8080.


export const UserCreate = {
   "payload": {
     "age": {
       "kind": "number",
       "minimum": 20,
       "maximum": 70
     },
     "email": {
       "kind": "string",
       "pattern": "^[a-z0-9._%+-]+@[a-z0-9.-]+.[a-z]{2,4}$"
     },
     "name": {
       "kind": "string",
       "min_length": 4,
       "max_length": 16
     },
     "sex": {
       "kind": "string",
       "enum": [
         "male",
         "female",
         "other"
       ]
     }
   }
 };

export const UserGet = {
   "UserID": {
     "kind": "number",
     "maximum": 10000
   }
 };

export const UserList = {};

export const RequiredError = "missing required parameter";
export const InvalidEnumValueError = "invalid enum value";
export const InvalidFormatError = "invalid format";
export const InvalidPatternError = "invalid pattern";
export const InvalidRangeError = "range exceeded";
export const InvalidMinLengthError = "length is less";
export const InvalidMaxLengthError = "length is exceeded";
export const InvalidKindError = "invalid kind";
export function validate(rule, actual) {
  let errors = {};

  if (typeof actual === "object") {
    Object.keys(actual).forEach(function(key, index) {
      const ret = validate(rule[key], actual[key]);
      if (ret !== undefined) {
        errors[key] = ret;
      }
    });
  } else {
    if (rule.kind !== typeof actual){
      errors.kind = InvalidKindError;
      return errors;
    }
    if (rule.maximum && actual > rule.maximum) {
      errors.maximum = InvalidRangeError;
    }
    if (rule.minimum && actual < rule.minimum) {
      errors.minimum = InvalidRangeError;
    }
    if (rule.max_length && actual.length > rule.max_length) {
      errors.max_length = InvalidMaxLengthError;
    }
    if (rule.min_length && actual.length < rule.min_length) {
      errors.min_length = InvalidMinLengthError;
    }
    if (rule.format && !(new RegExp(rule.format).test(actual))) {
      errors.format = InvalidFormatError;
    }
    if (rule.pattern && !(new RegExp(rule.pattern).test(actual))) {
      errors.pattern = InvalidPatternError;
    }
    if (rule.enum) {
      let found = false;
      for (let k in rule.enum) {
        if (k === actual) {found = true;}
      }
      if (found === false) {
        errors.enum = InvalidEnumValueError;
      }
    }
  }
  if (Object.keys(errors).length > 0){
    return errors;
  }
  return undefined;
}