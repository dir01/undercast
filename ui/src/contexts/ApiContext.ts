import * as preact from "preact";
import API from "../api";

const ApiContext = preact.createContext<API | undefined>(undefined);
export default ApiContext;
