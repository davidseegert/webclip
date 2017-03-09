# -*- coding: UTF-8 -*-

require 'FileUtils'
require 'nokogiri' # html parser
require 'liquid' # template language

class Template
    def initialize(name,variables=nil)
        # initialize variables
        @parse = Array.new()
        @copy = Array.new()
        @variables = Hash.new()
        @name = name

        templatePath = 'templates/'+name
        templateConfig =  JSON.parse(File.read(templatePath+'/config.json'))
        templateConfig['parse'].each do |folder|
            @parse.push(folder)
        end
        templateConfig['copy'].each do |folder|
            @copy.push(folder)
        end
        templateConfig['variables'].each do |key,variable|
            @variables[key] = variable
        end
        if variables != nil
            @variables = @variables.merge(variables)
        end
    end

    def setVariable(var,value)
      @variables[var] = value
    end

    def copy(projectFolder)
        # Copy folders from template
        puts 'Copy folders from template...'
        @copy.each do |folder|
          templatePath = File.join('templates',@name,folder)
          targetPath = File.join(projectFolder,'output',folder)
          FileUtils.mkdir_p targetPath
          FileUtils.copy_entry templatePath, targetPath
        end
        # Parsing files from template
        puts 'Begin parsing template files...'
        @parse.each do |folder|
            templatePath = File.join('templates',@name,folder)
            targetPath = File.join(projectFolder,'output',folder)
            FileUtils.mkdir_p targetPath
            #FileUtils.copy_entry File.join('templates',@name,folder), File.join(projectFolder,'output',folder)
            Dir.entries(templatePath).each do |file|
                templateFile = File.join(templatePath,file)
                targetFile = File.join(targetPath,file)
                if File.file?(templateFile)
                    system 'cls'
                    puts 'Parsing template file:'
                    puts targetFile
                    template = Liquid::Template.parse(File.read(templateFile))
                    rendered = template.render(@variables)

                    File.open(targetFile, 'w') { |wfile| wfile.write(rendered) }
                end
            end
        end


    end

    def parse(filePath)
        data = self.getData(filePath)
        self.getData(filePath)

        file = File.read('templates/'+@name+"/index.html")
        template = Liquid::Template.parse(file)
        rendered = template.render(data)
        return rendered

    end


    def getData(filePath)
        # data will contain all page-vaiables
        # adding template variables to data
        data = @variables.dup
        doc = Nokogiri::HTML(File.open(filePath))

        # get page title and body
        data['title'] = doc.xpath('//html//head//title').inner_html
        data['body'] = doc.xpath('//html//body').inner_html

        # get variables from meta tags
        metaTags = doc.xpath('//html//head//meta')
        metaTags.each do |meta|
            meta.keys.each do |key|
                if key[0..4] == 'data-'
                    value = meta[key]
                    if value.include? ";"
                        value = value.split ";"
                    end
                    data[key[5..-1]] = value
                end
            end
        end
        #puts data
        return data
    end

end
