/*
* This module is responsible to patch a modified file to include the work done by the 
* worker-dev. ex: how to merge two diffrent code/files that started of from one file, 
* and are currentely diffrent. diffrent due to work done by worker-dev and master-dev 
* in time elapsed
*/

var sfs = require('./SimpleFs.js');

/*
* _fn Makes Sure, the content is that of a text files and applicable for 
* patching.
* content is an array containing file content usualy read from disk.
*/
function test_if_text_content(content){
    for(var i = 0; i < content.length; i++){        
        if(content[i].split("\0").length > 1)
            return false;
    }
    return true;
}

function remove_indexes(main_array, indexes){
    if(indexes.length == 0)
        return main_array;
    var ret = [];
    for(var i = 0; i < main_array.length; i++)
        if(indexes.indexOf(i) == -1)
            ret.push(main_array[i]);
    return ret;
}

function test(o, d){
    var orignal_file = format_lines_of_code(sfs.read(o).toString());
    var developed_file = format_lines_of_code(sfs.read(d).toString());
    console.log(JSON.stringify(build_scope_patch(orignal_file, developed_file)));
}

function are_scopes_similar(o, d){
    if((o == undefined) && (d == undefined))
        return true;    
    if((o != undefined) && (d == undefined))
        return false;
    if((o == undefined) && (d != undefined))
        return false;    
    if(o.length != d.length)
        return false;    
    for(var i=0; i<o.length; i++){
        if(o[i].str != d[i].str)
            return false;
        if(!are_scopes_similar(o[i].block, d[i].block))
            return false;
    }
    return true;
}

/*
* _fn Builds a patch for new_file based on dev_file and orignal_file
* @param{string, file_path} d is the file path of the modified orignal file, mostly the work done 
* by the worker-dev.
* @param{string, file path} o is the file posted by the master-dev, the very first genesus file
*/
function build_patch(o, d, is_scope_tree, scope_name){
    var orignal_scope_tree = undefined;
    var developed_scope_tree = undefined;
    if(scope_name == undefined)
        scope_name = ["<libdevdev-global>"];
    if((is_scope_tree == undefined) || (is_scope_tree == false)){        
        var orignal_file = format_lines_of_code(sfs.read(o).toString());
        orignal_scope_tree = build_scope_tree(orignal_file);
        var developed_file = format_lines_of_code(sfs.read(d).toString());
        developed_scope_tree = build_scope_tree(developed_file);  
    }
    else{
        if(are_scopes_similar(o, d))
            return undefined;
        orignal_scope_tree = o;
        developed_scope_tree = d;
    }    
    var orignal_scope_1 = [];
    var developed_scope_1 = [];
    
    for(var i = 0 ; i < orignal_scope_tree.length; i++)
        orignal_scope_1.push(orignal_scope_tree[i].str);
    for(var i = 0; i < developed_scope_tree.length; i++)
        developed_scope_1.push(developed_scope_tree[i].str);
    
    var scope_1_patch = build_scope_patch(orignal_scope_1, developed_scope_1);
    
    var remove_from_developed_scope_tree = [];
    for(var i=0 ; i < scope_1_patch.main.length; i++){
        for(var j=scope_1_patch.main[i].code_index, k=0; 
            k < scope_1_patch.main[i].code.length; j++, k++)            
            if(developed_scope_tree[j].str == scope_1_patch.main[i].code[k]){
                if((developed_scope_tree[j].block != undefined) && 
                   (developed_scope_tree[j].block.length > 0)){
                    scope_1_patch.main[i].code[k] += merge_block(developed_scope_tree[j].block);
                }
                remove_from_developed_scope_tree.push(j);       
            }     
    }    
    developed_scope_tree = remove_indexes(developed_scope_tree, remove_from_developed_scope_tree);
    
    var remove_from_orignal_scope_tree = [];
    for(var i=0; i < scope_1_patch.help.length; i++){
        for(var j=0; j < scope_1_patch.help[i].code.length; j++) 
            if(orignal_scope_tree[scope_1_patch.help[i].code_index + j].str == scope_1_patch.help[i].code[j])
                remove_from_orignal_scope_tree.push(scope_1_patch.help[i].code_index + j);
    }
    orignal_scope_tree = remove_indexes(orignal_scope_tree, remove_from_orignal_scope_tree);
        
    //Now only similar scopes would be in the array of orignal_scope_tree and developed_scope_tree;    
    var j = 0;
    var scope_patches = [];
    //Push only if necessary.
    if((scope_1_patch != undefined) && 
       ((scope_1_patch.main.length > 0) || (scope_1_patch.help.length > 0)))
        scope_patches.push({name:scope_name, patch:scope_1_patch});
    
    for(var i=0; i < orignal_scope_tree.length; i++){
        if((orignal_scope_tree[i].block == undefined) || 
           (orignal_scope_tree[i].block.length == 0))
                continue;
        var temp_j = j;
        for(; j < developed_scope_tree.length; j++)
            if((developed_scope_tree[j].block != undefined) &&
               (developed_scope_tree[j].block.length != 0) &&
               (developed_scope_tree[j].str == orignal_scope_tree[i].str))
                break;
        if(j != developed_scope_tree.length){
            var scope_patch = build_patch(
                                orignal_scope_tree[i].block, 
                                developed_scope_tree[j].block, 
                                true, 
                                scope_name.concat([developed_scope_tree[j].str])
                            );
            //Accept only if got content
            if(scope_patch != undefined)
                for(var c = 0; c < scope_patch.length; c++)
                    if((scope_patch[c].patch.main.length > 0) || (scope_patch[c].patch.help.length > 0))
                        scope_patches.push(scope_patch[c]);
        }
        else
            j = temp_j;
    }
    return scope_patches;
}

exports.build_patch = build_patch;

function catTree(t){
    for(var i = 0; i < t.length; i++){
        console.log(t[i].str);
        if(t[i].block != undefined)
            catTree(t[i].block);
    }
}


function eq_array(arr1, arr2){
    if(arr1.length != arr2.length)
        return false;
    for(var i=0; i < arr1.length; i++){
        if(arr1[i] != arr2[i])
            return false;
    }
    return true;
}

function get_max_occurance_value(arr){
    var main_count = 0;
    var main_count_value = 0;
    for(var i=0; i<arr.length; i++){
        var count = 0;
        for(var j = i; j<arr.length; j++)
            if(arr[j] == arr[i])
                count++;
        if(count > main_count){
            main_count = count;
            main_count_value = arr[i];
        }
    }
    return main_count_value;
}

function get_index_of_min(arr, except){
    var min = Number.MAX_VALUE;
    var min_index = 0;
    for(var i = 0; i < arr.length; i++){
        if(arr[i] == except)
            continue;
        if(arr[i] <= min){
            min = arr[i]
            min_index = i;
        }
    }
    return min_index;
}

function get_above_negative_below_and_insert_index(scope, above, below){
    var all_above_index = [];
    var all_below_index = [];
    for(var i=0; i<scope.length; i++)
        if(scope[i].str == below)
            all_below_index.push(i);
        else if(scope[i].str == above)
            all_above_index.push(i);
    if((all_above_index.length == 0) && (all_below_index.length == 0))
        return [-1, (scope.length - 1)];
    if((all_above_index.length == 0) && (all_below_index.length > 0))
        return [Number.MAX_VALUE, all_below_index[0]];
    if((all_above_index.length > 0) && (all_below_index.length == 0))
        for(var i=0; i<all_above_index.length; i++)
            if(all_above_index[i] > 0)
                return [Number.MAX_VALUE, all_above_index[i] - 1];
    if((all_above_index.length > 0) && (all_below_index.length > 0)){
        var diff_len = Number.MAX_VALUE;
        var index = scope.length - 1; //Always Default to End of Scope..
        for(var a=0; a<all_above_index.length; a++)
            for(var b=0; b<all_below_index.length; b++)
                if(((all_above_index[a] - all_below_index[b]) <= diff_len) && 
                   ((all_above_index[a] - all_below_index[b]) > 0)){
                        diff_len = (all_above_index[a] - all_below_index[b]);
                        index = all_below_index[b];
                }
        return [diff_len, index];
    }
    return [-1, scope.length - 1];
}

//This function helps if one or more scopes have the same name.
function determine_right_scope_to_patch(possible_scopes, patch_code){
    if((possible_scopes == undefined) || (possible_scopes.length == 0))
        return undefined;
    if(possible_scopes.length == 1)
        return possible_scopes[0];
    var patch_code_possible_scopes = [];
    for(var i=0; i<patch_code.length; i++){
        var above = patch_code[i].above;
        var below = patch_code[i].below;        
        var best_fit_scope = [];        
        for(var j = 0; j<possible_scopes.length; j++){
            best_fit_scope.push(
                get_above_negative_below_and_insert_index(possible_scopes[j], above, below)[0]
            );
        }
        patch_code_possible_scopes.push(get_index_of_min(best_fit_scope, -1));                 
    }   
    //find the maximum occuring scope and return that.
    return possible_scopes[get_max_occurance_value(patch_code_possible_scopes)];
}


function find_scopes(name, scope_tree, name_index){
    if(scope_tree == undefined)
        return undefined;
    if((name.length <= 1) || (name.length <= name_index))
        return undefined;    
    var found = [];
    for(var i=0; i<scope_tree.length; i++){
        if(name[name_index] == scope_tree[i].str){
            if(name_index == (name.length - 1))
                found.push(scope_tree[i].block);
            else{
                var ret = find_scopes(name, scope_tree[i].block, name_index+1);
                if(ret != undefined)
                    found = found.concat(ret);
            }
        }    
    }
    return found;
}

/*
* _fn applies the patch to the new file recived from the master-dev. basicaly injecting
* work done by worker-dev in the new_file, preserving effort made by worker-dev.
* patch is the data structure representing the work done by worker-dev
* file_path is the new file that has to be modified;
*/
function apply_patch(patch, file_path){    
    var new_file_scope_tree = build_scope_tree(format_lines_of_code(
                                        sfs.read(file_path).toString())
                                    );
    var final_patch_code_and_scope = [];
    var remove_lns_from_scopes = []; 
    for(var i = 0; i< patch.length; i++){
        if((patch[i].name.length != 1) || (patch[i].name[0] != "<libdevdev-global>"))
            continue;
        remove_lns_from_scopes.push({_scope:new_file_scope_tree, _remove:patch[i].patch.help});
        for(var j = 0; j<patch[i].patch.main.length; j++){
            var above = patch[i].patch.main[j].above;
            var below = patch[i].patch.main[j].below;
            
            var insert_index_and_len = 
                get_above_negative_below_and_insert_index(new_file_scope_tree, above, below);
            var diff_len = 0; 
            if((insert_index_and_len[0] != Number.MAX_VALUE) && 
               (insert_index_and_len[0] != -1))
                diff_len = insert_index_and_len[0];
            
            final_patch_code_and_scope.push({
                _scope:new_file_scope_tree, 
                _in_between:diff_len,
                _index:insert_index_and_len[1],
                _code:patch[i].patch.main[j].code,
            });
        }                        
    }
    for(var i = 0; i < patch.length; i++){
        //We dont handel "libdevedev-global" patches hear, since its already done.
        if(patch[i].name.length <= 1)
            continue;           
        var scopes_to_patch = find_scopes(patch[i].name, new_file_scope_tree, 1);        
        var scope = determine_right_scope_to_patch(scopes_to_patch, patch[i].patch.main); 
        remove_lns_from_scopes.push({_scope:scope, _remove:patch[i].patch.help});
        for(var j=0; j<patch[i].patch.main.length; j++){
            var above = patch[i].patch.main[j].above;
            var below = patch[i].patch.main[j].below;
            var insert_index_and_len = get_above_negative_below_and_insert_index(scope, above, below);
            var diff_len = 0; 
            if((insert_index_and_len[0] != Number.MAX_VALUE) && 
               (insert_index_and_len[0] != -1))
                diff_len = insert_index_and_len[0];
            final_patch_code_and_scope.push({
                _scope:scope,
                _in_between:diff_len,
                _index:insert_index_and_len[1],
                _code:patch[i].patch.main[j].code,
            });
        }   
    }
    
    //First Remove the lines in _remove, from each scope.
    for(var i=0; i < remove_lns_from_scopes.length; i++){
        var scope = remove_lns_from_scopes[i]._scope;
        var remove_lns = remove_lns_from_scopes[i]._remove;
        for(var j=0; j<remove_lns.length; j++){
            var code = remove_lns[j].code;
            var code_index = remove_lns[j].code_index;
            for(var k=0; k<code.length; k++){
                var _add = 0;
                if(code_index >= scope.length)
                    _add += (code_index - (scope.length - 1));                
                var should_break = false;
                for(var s=0; s<scope.length+_add; s++){
                    if((code_index + s) < scope.length)
                        if(scope[code_index+s].str == code[k]){
                            scope.splice(code_index+s, 1); 
                            break;
                        }
                    if(((code_index - s) > -1) && ((code_index - s) < scope.length))
                        if(scope[code_index - s].str == code[k]){
                            scope.splice(code_index-s, 1);
                            break;
                        }
                }
            }            
        }
    }

    //Now Finaly We can start patching every single scope.
    for(var i=0; i < final_patch_code_and_scope.length; i++){
        var code = final_patch_code_and_scope[i]._code.join("\n");
        var scope = final_patch_code_and_scope[i]._scope;
        scope[final_patch_code_and_scope[i]._index].insert_after = code;
    }
    //There might be a extra '\n' character at index 0 added by merge blocks, remove that.
    var patched_file_produce = merge_block(new_file_scope_tree);
    patched_file_produce = patched_file_produce.substr(1, patched_file_produce.length);
    //Flush the new content to file..
    sfs.writeFile("./produce.txt", patched_file_produce);
}   
exports.apply_patch = apply_patch;


/*
* _fn gets the indentation count of the code line, ie # of occurences of '\t' or ' '
*/
function GetIndentationCount(line){
    if((line.length == 0) || (line.replace(/[\s\n]/g, '').length == 0))
        return undefined;
    var indentation_count = 0;
    for(var i = 0; i < line.length; i++)
        if((line.charAt(i) == ' '))
            indentation_count++;
        else if(line.charAt(i) == '\t')
            indentation_count += 4;
        else
            break;
    return indentation_count;
}
    
    
function merge_block(block){
    var str = "";
    for(var i=0; i < block.length; i++){
        str += "\n" + block[i].str
        if((block[i].block != undefined) && (block[i].block.length > 0))
            str += merge_block(block[i].block); 
        if((block[i].insert_after != undefined) && (block[i].insert_after != null))
            str += "\n"+block[i].insert_after;
    }
    return str;
}


/*
* _fn build the indentation tree, 
* ie datastructures that represent indentation in file, since almost all code and 
* project files have indentation
* linesofcode is a array of strings containing files content
*/
function build_scope_tree(linesofcode){
    if(linesofcode.length == 0)
        return undefined;
    var blocks = [];
    var block_linesofcode = [];    
    var main_indentation_count = undefined;
    
    var i = 0;
    for(i = 0; i < linesofcode.length; i++){
        blocks.push({str:linesofcode[i], block:undefined});
        main_indentation_count = GetIndentationCount(linesofcode[i]);
        if(main_indentation_count != undefined)
            break;
    }    
    for(i += 1; i < linesofcode.length; i++){
        var indentation_count = GetIndentationCount(linesofcode[i]);
        if((indentation_count != undefined) && (indentation_count <= main_indentation_count)){
            if(blocks[blocks.length - 1] != undefined)
                blocks[blocks.length - 1].block = build_scope_tree(block_linesofcode);
            blocks.push({str:linesofcode[i], block:undefined});
            block_linesofcode = [];
            main_indentation_count = indentation_count;
        }
        else
            block_linesofcode.push(linesofcode[i]);    
    }    
    if(blocks[blocks.length - 1] != undefined)
        blocks[blocks.length - 1].block = build_scope_tree(block_linesofcode);
    return blocks;
}


function format_lines_of_code(str){
    var linesofcode = str.split('\n');
    for(var i = 0; i < linesofcode.length; i++){
        if((i  > 0) && ((i + 1) < linesofcode.length)){
            var current_indentation = GetIndentationCount(linesofcode[i]);
            var next_indentation = GetIndentationCount(linesofcode[i+1]);
            if( (linesofcode[i].search(/.*[0-9a-zA-Z].*/) == -1) && 
                (next_indentation > current_indentation) && 
                (linesofcode[i].trim() != "")
              ){
                linesofcode[i - 1] = linesofcode[i - 1] + "\n" + linesofcode.splice(i, 1)[0];
                i--;
            }    
        }
    }
    return linesofcode;
}

function build_scope_patch(o, d){
    var developed_file = d;
    var orignal_file = o;
    var common_code = lcs(orignal_file, developed_file);
    var delta = [];
    var delta_2 = [];
    var temp_i = -2;
    for(var i = 0, j = 0; i < developed_file.length; i++)
        if(developed_file[i] == common_code[j])
            j++
        else{
            if((i - temp_i) > 1){
                if((delta[delta.length - 1] != undefined) && (temp_i > 0) && (temp_i < (developed_file.length - 1)))
                    delta[delta.length - 1].above = developed_file[temp_i + 1];     
                delta.push({code:[], code_index:i});
                if(i > 0)
                    delta[delta.length - 1].below = developed_file[i - 1];         
            }            
            delta[delta.length - 1].code.push(developed_file[i]);
            temp_i = i;
        }
    if((delta[delta.length - 1] != undefined) && (temp_i > 0) && (temp_i < (developed_file.length - 1)))
        delta[delta.length - 1].above = developed_file[temp_i + 1]; 
    
    //Build the Helper
    temp_i = -2;
    for(var i = 0, j = 0; i < orignal_file.length; i++)
        if(orignal_file[i] == common_code[j])
            j++
        else{
            if((i - temp_i) > 1)
                delta_2.push({code:[], code_index:i});         
            delta_2[delta_2.length - 1].code.push(orignal_file[i]);
            temp_i = i;
        }
    
    return {main:delta, help:delta_2};
}

/*
* _fn builds an array of longest common series. present in x and y
* x & y are arrays
*/
function lcs(x,y){
	var s,i,j,m,n,
		lcs=[],row=[],c=[],
		left,diag,latch;
	//make sure shorter string is the column string
	if(m<n){s=x;x=y;y=s;}
	m = x.length;
	n = y.length;
	//build the c-table
	for(j=0;j<n;row[j++]=0);
	for(i=0;i<m;i++){
		c[i] = row = row.slice();
		for(diag=0,j=0;j<n;j++,diag=latch){
			latch=row[j];
			if(x[i] == y[j]){row[j] = diag+1;}
			else{
				left = row[j-1]||0;
				if(left>row[j]){row[j] = left;}
			}
		}
	}
	i--,j--;
	//row[j] now contains the length of the lcs
	//recover the lcs from the table
	while(i>-1&&j>-1){
		switch(c[i][j]){
			default: j--;
				lcs.unshift(x[i]);
			case (i&&c[i-1][j]): i--;
				continue;
			case (j&&c[i][j-1]): j--;
		}
	}
	return lcs;
}
