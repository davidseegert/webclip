# -*- coding: UTF-8 -*-

require_relative 'Template'

class Project
    def initialize(path)
        # parse json file
        configFilePath = File.join(path,'config.json')
        configFile = File.read(configFilePath, :encoding => 'UTF-8')
        begin
          config = JSON.parse(configFile)
        rescue
          printe("Der Projekt-Ordner \""+path+"\" enthält eine ungültige config.json Datei")
        end
        # create new projet-object and add it to projects-array
        @config = config
        @name = config['name']
        @template = config['template']
        @variables = config['variables']
        @rewrites = config['rewrites']
        @path = path
        @parseFiles = Array.new()
        @menus = {
            'main' => Array.new()
        }
        @outputPath = 'projects/'+@name+"/output/"

    end

    def getName
        return @name
    end

    def getPath
        return @path
    end

    def getMenus
        @parseFiles.each do |file|
            #puts file
            doc = Nokogiri::HTML(File.open(file))


            # get variables from meta tags
            menuTags = doc.xpath('//html//head//meta/@data-menu')

            # get variables from meta tags
            titles = doc.xpath('//html//head//title')
            if titles.length > 0
            title = titles[0].inner_html
            else
              title = nil
            end

            if title == nil
                title = '! '+' [h1] fehlt'
            end

            menuWeights = doc.xpath('//html//head//meta/@data-weight')
            #puts menuWeights.length
            if menuWeights.length > 0
              weight = menuWeights[0].value.to_i
            else
              weight = nil
            end



            if weight == nil
              weight = 0
            end

            sort = nil
            data_sort = doc.xpath('//html//head//meta/@data-sort')
            if data_sort.length > 0
              sort = data_sort[0].value
            end

            e = {
              'title' =>  title,
              'href' => file[@outputPath.length..-1],
              'weight' => weight,
              'sort' => sort
            }


            #puts e

            menuTags.each do |menuTag|
                # create key if nox exists
                if @menus.key?(menuTag.value) == false
                    @menus[menuTag.value] = Array.new()
                end
                @menus[menuTag.value].push(e)
            end
            if menuTags.length == 0
                @menus['main'].push(e)
            end

            # adding aliases
            menuAlias = doc.xpath('//html//head//meta/@data-alias')
            #puts menuWeights.length
            menuAlias.each do |entrie|
              ealias = e.clone
              ealias['alias'] = entrie.value
              @menus['main'].push(ealias)
            end

        end
    end

    def renderMenu(name,className,current)
      # return empty array if menu not exists
      if @menus[name] == nil
        return Array.new()
      end

      # order menu entries by weight attribute
      entries = @menus[name]
      sorted = entries.sort_by{|k| k['weight']}

      # create html menu from array
      out = '<ul class="'+className+'">'
      sorted.each{|entry|
        active = ''
        if entry['href'] == current
          active = ' class="active" '
        end
        out += '<li><a '+active+' href="'+entry['href']+'">'+entry['title']+'</a></li>'
      }
      out += '</ul>'
      return out
    end

    def renderMenuTopic(name,className,current)
      begin
        topics = @variables['topics'].clone
      rescue
        topics = Array.new()
      end
      topics.push("???")
      allEntries = @menus[name]
      allEntries.each do |e|
        if e['sort'] == nil
          e['sort'] = "???"
        end
      end
      out = ''

      topics.each do |topic|
        unsortedEntries = Array.new()
        allEntries.each do |e|
          if e['sort'] == topic
            unsortedEntries.push(e);
          end
        end

        if topic == "???" and unsortedEntries.count == 0
          next
        end

        sortedEntries = unsortedEntries.sort_by { |k| k['weight'] }
        out += "<h3>#{topic}</h3><ul>"
        sortedEntries.each do |e|
          active = ''
          if e['href'] == current
            active = ' class="active" '
          end
          out += "<li><a "+active+"href=\"#{e['href']}\">#{e['title']}</a></li>"
        end
        out += "</ul>"
      end

      return out
    end

    def renderMenuAlphabetical(name,className,current)
        #entries = @menus[name].sort_by { |word| word.downcase }
        entries = @menus[name]
        #return entries
        sorted = Array.new()
        entries.each{|entry|
              # removing umlaute for sorting in navigation
              sort = entry['title'].encode('UTF-8')
              sort = sort.gsub('Ä', 'Ae')
              sort = sort.gsub('ä', 'ae')
              sort = sort.gsub('Ö', 'Oe')
              sort = sort.gsub('ö', 'oe')
              sort = sort.gsub('Ü', 'Ue')
              sort = sort.gsub('ü', 'ue')
              sort = sort.gsub('²', '2')
              sort = sort.gsub(/[^0-9A-Za-z]/, '')

              entry['sort'] = sort

              # Appent Alias to sort Name
              if(entry.has_key? "alias")
                aliasSort = entry['alias']
                aliasSort = aliasSort.gsub('Ä', 'Ae')
                aliasSort = aliasSort.gsub('ä', 'ae')
                aliasSort = aliasSort.gsub('Ö', 'Oe')
                aliasSort = aliasSort.gsub('ö', 'oe')
                aliasSort = aliasSort.gsub('Ü', 'Ue')
                aliasSort = aliasSort.gsub('ü', 'ue')
                aliasSort = aliasSort.gsub('²', '2')
                sort = sort.gsub('²', '2')
                aliasSort = aliasSort.gsub(/[^0-9A-Za-z]/, '')


                entry['sort'] = aliasSort
              end
        }
        sorted = entries.sort_by{|k| k['sort'].downcase }

        letter = ''
        out = '<ul>'
        for item in sorted
            if(letter != item['sort'][0].upcase)
                letter = item['sort'][0].upcase
                out += '</ul>'
                out += '<h2>'+letter+'</h2>'
                out += '<ul>'
            end
            active = ''
            if item['href'] == current
              active = ' class="active" '
            end
            ali = ""
            if(item.has_key? "alias")
              ali = '<span style="color:gray">'+item['alias']+' ⟼ </span>'
            end
            out += "<li>#{ali}<a "+active+"href=\"#{item['href']}\">#{item['title']}</a></li>"
        end
        out += "</ul>"
        return out
    end

    def createSearch()
      data = {}
      @parseFiles.each do |file|
          system 'cls'
          puts "Erstelle Suche..."
          puts file
          doc = Nokogiri::HTML(File.open(file))
          filename = file[@outputPath.length..-1]

          # get variables from meta tags
          tags = doc.xpath('//meta[@name="keywords"]/@content')
          begin
            title = doc.xpath('//title').first.inner_html
          rescue
            printe("Es scheint als hätte die Datei \"#{file}\" kein <title>-Tag.");
          end
          aliases = doc.xpath('//html//head//meta/@data-alias')


          data[filename] = Hash.new()
          data[filename]['tags'] = Array.new()
          data[filename]['title'] = title
          data[filename]['aliases'] = Array.new()

          aliases.each do |al|
            data[filename]['aliases'].push(al.value)
          end
          tags.each do |tag|
            tagArray = tag.value().split(',')
            tagArray.each do |singleTag|
              data[filename]['tags'].push(singleTag)
            end

          end
      end
      File.open(@outputPath+"search.json", 'w') { |file| file.write(data.to_json) }
    end

    def render

        # copy raw files to output folder
        puts 'Deleting old files...'
        FileUtils.rm_rf(@outputPath)
        puts 'Copy source files...'
        FileUtils.copy_entry @path+'/src', @outputPath




        # get all files to be parsed
        @config['parse'].each do |files|
            Dir.glob(@path+"/output/"+files).each{ |file|
                @parseFiles.push(file)
            }
        end

        # template selection and variable passing
        template = Template.new('default',@variables)
        # copy template files to project output
        template.copy(@path)

        self.getMenus
        # generating the search.json file
        createSearch()

        #parsing of template files
        @parseFiles.each do |file|
            system 'cls'
            puts 'Erstelle Datei:'
            puts file
            if @variables.has_key?('topics') == false
              template.setVariable('mainmenu',self.renderMenuAlphabetical('main','',file[@outputPath.length..-1]))
            else
              template.setVariable('mainmenu',self.renderMenuTopic('main','',file[@outputPath.length..-1]))
            end
            template.setVariable('topmenu',self.renderMenu('top','topmenu',file[@outputPath.length..-1]))
            parsed = template.parse(file)
            #FileUtils.rm_rf(outPath)
            File.open(file, 'w') { |wfile| wfile.write(parsed) }

        end

        return @outputPath
    end

end
