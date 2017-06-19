
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

interface UserMedia {
  email: string;
  age: number;
  sex: string;
  name: string;
}

interface UserTypeCollectionMedia {
  name: string;
  email: string;
  age: number;
  sex: string;
}
