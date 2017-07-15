
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
  name: string;
  age: number;
  email: string;
  sex: ["male","female","other"];
}

interface UserMedia {
  sex: string;
  name: string;
  email: string;
  age: number;
}

interface UserTypeCollectionMedia {
  name: string;
  email: string;
  age: number;
  sex: string;
}
