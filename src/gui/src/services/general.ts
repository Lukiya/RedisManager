import u from '@/u';
import { request } from 'umi';

export async function GetInfo() {
    const r = await request(u.LocalRootURL());

    return r;
}