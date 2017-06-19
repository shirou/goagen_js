
declare class ErrorMap {
   kind?: string;
   maximum?: string;
   minimum?: string;
   max_length?: string;
   min_length?: string;
   format?: string;
   pattern?: string;
   enum?: string;
   [key: string]: ErrorMap | string | undefined;
}

interface UserCreatePayload {
  sex: ["male","female","other"];
  name: string;
  age: number;
  email: string;
}
