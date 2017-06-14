import React from 'react';
import {Field, reduxForm} from 'redux-form';

import * as v from "../api_validator.ts";

const renderField = ({input, label, type, meta: {touched, error, warning}}) => (
  <div>
    <label>{label}</label>
    <div>
      <input {...input} placeholder={label} type={type} />
      {touched &&
        ((error && <span>{error}</span>) ||
          (warning && <span>{warning}</span>))}
    </div>
  </div>
);

const v_wrap = (rule: any, value: any): string | undefined => {
  const e = v.validate(rule, value);
  if (e === undefined){ return undefined };
  if (e.required){
    return "required";
  }
  if (e.maximum){
    return "too old";
  }
  if (e.minimum){
    return "too young";
  }
  if (e.min_length){
    return "should be more than 4";
  }
  if (e.max_length){
    return "should be less than 16";
  }
  if (e.kind){
    return e.kind;
  }
  if (e.pattern) {
    return "invalid format string";
  }
};

const name = (value: string) => v_wrap(v.UserCreate.payload.name, value);
const email = (value: string) => v_wrap(v.UserCreate.payload.email, value);
const age = (value: string) => v_wrap(v.UserCreate.payload.age, parseInt(value, 10));

const FieldLevelValidationForm = props => {
  const {handleSubmit, pristine, reset, submitting} = props;
  return (
    <form onSubmit={handleSubmit}>
      <Field
        name="name"
        type="text"
        component={renderField}
        label="Name"
        validate={name}
      />
      <Field
        name="email"
        type="email"
        component={renderField}
        label="Email"
        validate={email}
        warn={email}
      />
      <Field
        name="age"
        type="number"
        component={renderField}
        label="Age"
        validate={age}
        normalize={value => parseInt(value, 10)}
      />
      <div>
        <button type="submit" disabled={submitting}>Submit</button>
        <button type="button" disabled={pristine || submitting} onClick={reset}>
          Clear Values
        </button>
      </div>
    </form>
  );
};

export default reduxForm({
  form: 'fieldLevelValidation' // a unique identifier for this form
})(FieldLevelValidationForm);
