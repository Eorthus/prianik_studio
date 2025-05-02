import { computed, ref } from "vue";

const categories = ref();

const categoriesMap = computed(() =>
  categories.value?.map((el) => {
    return { id: el.id, name: el.name };
  })
);

const subcategoriesMap = computed(() => {
    let subs = []
    categories.value?.forEach((el) => {
        if (el.subcategories?.length) {
            subs.push(el.subcategories.map(sub=>{
                return { id: sub.id, name: sub.name, parent_id: sub.parent_id };
            }))
        }
      })
      return subs.flat()
}
);

const getCategory = (id?:number)=>{
    if (!id) {
        return ''
    }

    return categoriesMap.value?.find(el=>el.id===id)?.name
}

const getSubcategory = (id?:number, parent_id?:number)=>{
    if (!id || !parent_id) {
        return ''
    }
    return subcategoriesMap.value?.find(el=>el.id===id && el.parent_id===parent_id)?.name
}

export const useCategories = () => {
  return {
    categories,
    categoriesMap,
    subcategoriesMap,
    getCategory,
    getSubcategory,
  };
};
