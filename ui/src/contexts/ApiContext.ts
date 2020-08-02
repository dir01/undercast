import * as preact from "preact";
import API from "../API";

const ApiContext = preact.createContext<API | undefined>(undefined);
export default ApiContext;
