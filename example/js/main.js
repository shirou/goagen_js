import React from 'react';
import ReactDOM from 'react-dom';
import {Provider} from 'react-redux';
import {createStore, combineReducers} from 'redux';
import {reducer as reduxFormReducer} from 'redux-form';
import { Values } from 'redux-form-website-template';

const dest = document.getElementById('content');
const reducer = combineReducers({
  form: reduxFormReducer // mounted under "form"
});
const store = (createStore)(reducer);

const showResults = values =>
  new Promise(resolve => {
    setTimeout(() => {
      // simulate server latency
      window.alert(`You submitted:\n\n${JSON.stringify(values, null, 2)}`)
      resolve()
    }, 500)
  });

const FieldLevelValidationForm = require('./components/FieldLevelValidationForm').default;

ReactDOM.render(
  <Provider store={store}>
    <div>
      <h2>Form</h2>
      <FieldLevelValidationForm onSubmit={showResults} />
      <Values form="fieldLevelValidation" />
    </div>
  </Provider>,
  dest
);
