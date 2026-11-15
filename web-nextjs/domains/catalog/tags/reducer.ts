import {Action, createGenericReducer} from "@/libs/daction";
import {TagPageState} from "./language";

export type TagPageAction = Action<TagPageState, any>

export const tagPageReducer = createGenericReducer<TagPageState>();
