import * as api from "./api_request.js";
import * as v from "./api_validator.js";

const payload = {
//  int_required: 10,
  int_enum: 10,
  int_max: 1,
};

//console.log(v.validate(v.GetGetInt.payload, payload));

console.log(v.validate(v.GetPathParams.ParamInt, 10000));
/*

api.GetGetInt(payload).then((response) => {
  if (response.status !== 200){
    console.log("error", response.status);
    throw response.json();
  }
  return {};
});
*/
